package posts

import (
	"database/sql"
	"log/slog"

	"go.temporal.io/sdk/workflow"
)

type WorkflowManager struct {
	db *sql.DB
}

func NewWorkflowManager(db *sql.DB) *WorkflowManager {
	return &WorkflowManager{db}
}

func (m *WorkflowManager) SendPostToTopicFollowers(ctx workflow.Context, post_id int, topics []int) error {
	slog.Info("Executing SendPostToTopicFollowers")
	// TODO: Get store, get followers for all the topics and send the post
	return nil
}
