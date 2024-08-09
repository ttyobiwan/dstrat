package posts

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ttyobiwan/dstrat/internal/sqlite"
)

var ErrNotFound = errors.New("entry not found")

type TopicStore interface {
	Create(name string) (*Topic, error)
	GetByName(name string) (*Topic, error)
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
	query := `SELECT id, name FROM topics WHERE name = ?`

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
