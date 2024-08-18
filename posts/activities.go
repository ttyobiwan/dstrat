package posts

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/ttyobiwan/dstrat/users"
)

type ActivityManager struct {
	db *sql.DB
}

func NewActivityManager(db *sql.DB) *ActivityManager {
	return &ActivityManager{db}
}

func (m *ActivityManager) TopicStore() TopicStore {
	return NewTopicDBStore(m.db)
}

func (m *ActivityManager) PostStore() PostStore {
	return NewPostDBStore(m.db)
}

func (m *ActivityManager) GetPost(id int) (*Post, error) {
	slog.Info("Getting post", "id", id)

	post, err := m.PostStore().Get(id)
	if err != nil {
		return nil, fmt.Errorf("getting post: %v", err)
	}
	return post, nil
}

func (m *ActivityManager) GetFollowers(topics []int) ([]*users.User, error) {
	slog.Info("Getting followers", "topics", topics)

	followers, err := m.TopicStore().GetFollowers(topics)
	if err != nil {
		return nil, fmt.Errorf("getting followers: %v", err)
	}
	return followers, nil
}

func (m *ActivityManager) SendSinglePost(post *Post, follower *users.User) error {
	slog.Info("Sending single post", "post_id", post.ID, "follower_id", follower.ID)
	return SendSinglePost(post, follower)
}

func (m *ActivityManager) SendPostSequentially(post *Post, followers []*users.User) error {
	slog.Info("Sending post sequentially", "post_id", post.ID, "followers_count", len(followers))

	for _, follower := range followers {
		if err := SendSinglePost(post, follower); err != nil {
			slog.Error("Error sending post", "error", err)
		}
	}

	return nil
}

func (m *ActivityManager) SendPostASequentially(post *Post, followers []*users.User) error {
	slog.Info("Sending post asequentially", "post_id", post.ID, "followers_count", len(followers))

	// Very naive implementation - should at least have max parameter
	for _, follower := range followers {
		go func(follower *users.User) {
			if err := SendSinglePost(post, follower); err != nil {
				slog.Error("Error sending post", "error", err)
			}
		}(follower)
	}

	return nil
}
