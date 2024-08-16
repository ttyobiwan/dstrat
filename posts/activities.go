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

// TODO:
//   - Send single
//   - Send batch
//   - Send all
func (m *ActivityManager) DispatchSendPost(post *Post, followers []*users.User) error {
	slog.Info("Dispatching send post", "post_id", post.ID, "followers_count", len(followers))
	return nil
}
