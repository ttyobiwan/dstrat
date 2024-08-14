package posts

import (
	"context"

	"go.temporal.io/sdk/workflow"
)

type TaskClient[options any, result any] interface {
	Execute(ctx context.Context, task any, options options, args ...any) (result, error)
	AwaitResult(ctx context.Context, task result, dest any) error
	Close() error
}

func SendPostToTopicFollowers(ctx workflow.Context, post_id int, topics []int) error {
	return (&WorkflowManager{}).SendPostToTopicFollowers(ctx, post_id, topics)
}
