package users

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/ttyobiwan/dstrat/internal/sqlite"
)

var ErrNotFound = errors.New("entry not found")

type UserStore interface {
	Create(username string) (*User, error)
	GetByUsername(username string) (*User, error)
}

type UserDBStore struct {
	*sqlite.DBStore
}

func NewUserDBStore(db *sql.DB) *UserDBStore {
	return &UserDBStore{sqlite.NewDBStore(db)}
}

func (s *UserDBStore) Create(username string) (*User, error) {
	stmt, err := s.DB().Prepare(`INSERT INTO users (username) VALUES (?)`)
	if err != nil {
		return nil, fmt.Errorf("preparing statement: %v", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(username)
	if err != nil {
		return nil, fmt.Errorf("executing query: %v", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("getting last id: %v", err)
	}

	return &User{int(id), username}, nil
}

func (s *UserDBStore) GetByUsername(username string) (*User, error) {
	query := `SELECT id, username FROM users WHERE username = ?`

	var user User
	err := s.DB().
		QueryRow(query, username).
		Scan(&user.ID, &user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("scanning row: %v", err)
	}

	return &user, nil
}
