package posts

import (
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/ttyobiwan/dstrat/users"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type WorkflowManager struct {
	db         *sql.DB
	activities *ActivityManager
}

func NewWorkflowManager(db *sql.DB) *WorkflowManager {
	return &WorkflowManager{db, &ActivityManager{}}
}

func (m *WorkflowManager) TopicStore() TopicStore {
	return NewTopicDBStore(m.db)
}

func (m *WorkflowManager) PostStore() PostStore {
	return NewPostDBStore(m.db)
}

func (m *WorkflowManager) DefaultContext(ctx workflow.Context) workflow.Context {
	retrypolicy := &temporal.RetryPolicy{
		InitialInterval:    time.Second,
		BackoffCoefficient: 2.0,
		MaximumInterval:    1 * time.Minute,
		MaximumAttempts:    3,
	}
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy:         retrypolicy,
	}
	return workflow.WithActivityOptions(ctx, options)
}

func (m *WorkflowManager) SendPostToTopicFollowers(ctx workflow.Context, postID int, topics []int) error {
	slog.Info("Sending post to followers", "post_id", postID, "topics", topics)

	ctx = m.DefaultContext(ctx)

	// Get post
	post := Post{}
	workflow.ExecuteActivity(ctx, m.activities.GetPost, postID).Get(ctx, &post)

	// Get followers
	followers := []*users.User{}
	workflow.ExecuteActivity(ctx, m.activities.GetFollowers, topics).Get(ctx, &followers)

	// TODO: Move it somewhere
	strat := SendPostStratMass

	// Dispatch sending post
	switch strat {
	case SendPostStratSeq:
		workflow.ExecuteActivity(ctx, m.activities.SendPostSequentially, post, followers).Get(ctx, nil)
	case SendPostStratASeq:
		workflow.ExecuteActivity(ctx, m.activities.SendPostASequentially, post, followers).Get(ctx, nil)
	case SendPostStratBulk:
		workflow.ExecuteChildWorkflow(ctx, m.SendPostBulk, post, followers).Get(ctx, nil)
	case SendPostStratMass:
		workflow.ExecuteChildWorkflow(ctx, m.SendPostMass, post, followers).Get(ctx, nil)
	default:
		return fmt.Errorf("invalid sending strategy: %s", strat)
	}

	slog.Info("Done sending post to followers", "post_id", post.ID)

	return nil
}

func (m *WorkflowManager) SendPostBulk(ctx workflow.Context, post *Post, followers []*users.User) error {
	slog.Info("Sending post bulk", "post_id", post.ID, "followers_count", len(followers))
	// TODO: Implement this one
	return nil
}

func (m *WorkflowManager) SendPostMass(ctx workflow.Context, post *Post, followers []*users.User) error {
	slog.Info("Sending post mass", "post_id", post.ID, "followers_count", len(followers))

	ctx = m.DefaultContext(ctx)
	results := make([]struct{}, 0, len(followers))

	for _, follower := range followers {
		workflow.Go(ctx, func(gCtx workflow.Context) {
			workflow.ExecuteActivity(gCtx, m.activities.SendSinglePost, post, follower).Get(gCtx, nil)
			results = append(results, struct{}{})
		})
	}

	_ = workflow.Await(ctx, func() bool {
		return len(results) == len(followers)
	})

	slog.Info("Done sending post using mass strategy", "post_id", post.ID)
	return nil
}
