package posts

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/ttyobiwan/dstrat/internal/sqlite"
	"github.com/ttyobiwan/dstrat/users"
)

var ErrNotFound = errors.New("entry not found")

type TopicStore interface {
	Create(name string) (*Topic, error)
	GetByName(name string) (*Topic, error)
	ToggleFollow(topic_id string, user_id string) error
}

type TopicDBStore struct {
	*sqlite.DBStore
}

func NewTopicDBStore(db *sql.DB) *TopicDBStore {
	return &TopicDBStore{sqlite.NewDBStore(db)}
}

func (s *TopicDBStore) Create(name string) (*Topic, error) {
	stmt, err := s.DB().Prepare(`INSERT INTO topics (name) VALUES (?)`)
	if err != nil {
		return nil, fmt.Errorf("preparing statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(name)
	if err != nil {
		return nil, fmt.Errorf("executing query: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting last id: %v", err)
	}

	return &Topic{int(id), name}, nil
}

func (s *TopicDBStore) GetByName(name string) (*Topic, error) {
	const query = `SELECT id, name FROM topics WHERE name = ?`

	var topic Topic
	err := s.DB().
		QueryRow(query, name).
		Scan(&topic.ID, &topic.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("scanning row: %v", err)
	}

	return &topic, nil
}

func (s *TopicDBStore) ToggleFollow(topic_id string, user_id string) error {
	var rowid int
	err := s.
		DB().
		QueryRow(`SELECT rowid FROM topic_followers WHERE topic_id = ? AND user_id = ?`, topic_id, user_id).
		Scan(&rowid)
	// If there is an error, that is not ErrNoRows, then return it
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("scanning topic follow: %v", err)
		}
	}
	// Otherwise, delete the current follow
	if rowid != 0 {
		_, err = s.DB().Exec(`DELETE FROM topic_followers WHERE topic_id = ? AND user_id = ?`, topic_id, user_id)
		if err != nil {
			return fmt.Errorf("deleting topic follow: %v", err)
		}
		return nil
	}

	// Create a new follow
	stmt, err := s.DB().Prepare(`INSERT INTO topic_followers (topic_id, user_id) VALUES (?, ?)`)
	if err != nil {
		return fmt.Errorf("preparing statement: %v", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(topic_id, user_id)
	if err != nil {
		fmt.Println(err)
		return fmt.Errorf("executing query: %v", err)
	}

	return nil
}

type PostStore interface {
	Create(title, content string, author int, topics []int) (*Post, error)
}

type PostDBStore struct {
	*sqlite.DBStore
}

func NewPostDBStore(db *sql.DB) *PostDBStore {
	return &PostDBStore{sqlite.NewDBStore(db)}
}

func (s *PostDBStore) Create(title, content string, author int, topics []int) (*Post, error) {
	// Start the transaction
	err := s.BeginTx(context.Background())
	if err != nil {
		return nil, fmt.Errorf("beginning transaction: %v", err)
	}
	defer s.Rollback()

	// First, create the post object
	stmt, err := s.DB().Prepare(`INSERT INTO posts (title, content, author_id) VALUES (?, ?, ?)`)
	if err != nil {
		return nil, fmt.Errorf("preparing post statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(title, content, author)
	if err != nil {
		return nil, fmt.Errorf("executing post query: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting last id: %v", err)
	}

	// Now, insert into post-topic many-to-many
	postTopicsQuery := "INSERT INTO post_topics(post_id, topic_id) VALUES "
	var postTopicsInsert []string
	var postTopicsValues []any

	const rowSQL = "(?, ?)"
	for _, topicID := range topics {
		postTopicsInsert = append(postTopicsInsert, rowSQL)
		postTopicsValues = append(postTopicsValues, id, topicID)
	}
	postTopicsQuery = postTopicsQuery + strings.Join(postTopicsInsert, ",")

	stmt, err = s.DB().Prepare(postTopicsQuery)
	if err != nil {
		return nil, fmt.Errorf("preparing post-topics statement: %v", err)
	}

	_, err = stmt.Exec(postTopicsValues...)
	if err != nil {
		return nil, fmt.Errorf("executing post-topics query: %v", err)
	}

	// Finally, get full post object
	post, err := s.Get(int(id))
	if err != nil {
		return nil, fmt.Errorf("getting post: %v", err)
	}

	// Commit the transaction
	err = s.Commit()
	if err != nil {
		return nil, fmt.Errorf("commiting transaction: %v", err)
	}

	return post, err
}

func (s *PostDBStore) Get(id int) (*Post, error) {
	author := users.User{}
	post := Post{Author: &author}
	topics := ""

	err := s.DB().
		QueryRow(
			`SELECT
				p.id,
				p.title,
				p.content,
				u.id,
				u.username,
				GROUP_CONCAT(t.id || ':' || t.name, ';') AS topics
			FROM
				posts p
				LEFT JOIN users u ON p.author_id = u.id
				LEFT JOIN post_topics pt ON p.id = pt.post_id
				LEFT JOIN topics t ON pt.topic_id = t.id
			WHERE
				p.id = ?
			GROUP BY
				p.id`,
			id,
		).
		Scan(&post.ID, &post.Title, &post.Content, &author.ID, &author.Username, &topics)
	if err != nil {
		return nil, fmt.Errorf("scanning post: %v", err)
	}

	// Parse topics
	post.Topics = sqlite.ParseGroupConcat(topics, ";", ":", func(v []string) *Topic {
		tid, _ := strconv.Atoi(v[0])
		return &Topic{tid, v[1]}
	})

	return &post, nil
}
