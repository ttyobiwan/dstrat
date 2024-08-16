package posts

import (
	"database/sql"
	"log/slog"
	"time"

	"github.com/ttyobiwan/dstrat/users"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

type WorkflowManager struct {
	db *sql.DB
}

func NewWorkflowManager(db *sql.DB) *WorkflowManager {
	return &WorkflowManager{db}
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

	activities := &ActivityManager{}
	ctx = m.DefaultContext(ctx)

	// Get post
	post := Post{}
	workflow.ExecuteActivity(ctx, activities.GetPost, postID).Get(ctx, &post)

	// Get followers
	followers := []*users.User{}
	workflow.ExecuteActivity(ctx, activities.GetFollowers, topics).Get(ctx, &followers)

	// Dispatch send post
	workflow.ExecuteActivity(ctx, activities.DispatchSendPost, post, followers).Get(ctx, nil)

	return nil
}
