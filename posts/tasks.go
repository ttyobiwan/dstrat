package posts

import (
	"context"
	"errors"
	"log/slog"
	"math/rand/v2"
	"time"

	"github.com/ttyobiwan/dstrat/users"
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

func SendSinglePost(post *Post, follower *users.User) error {
	slog.Info("Sending post", "post_id", post.ID, "follower_id", follower.ID)

	time.Sleep(3 * time.Second)

	if rand.IntN(2) == 1 {
		return errors.New("post could not be sent")
	}

	slog.Info("Post successfully sent", "post_id", post.ID, "follower_id", follower.ID)

	return nil
}
