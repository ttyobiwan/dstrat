package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ConnOrTx is an interface that represents either db connection or transaction.
type ConnOrTx interface {
	Query(string, ...any) (*sql.Rows, error)
	QueryRow(string, ...any) *sql.Row
	Prepare(string) (*sql.Stmt, error)
}

// DBStore is a generic store (data access layer) that gives additional functionality of using transactions.
type DBStore struct {
	db *sql.DB
	tx *sql.Tx
}

func NewDBStore(db *sql.DB) *DBStore {
	return &DBStore{db, nil}
}

// DB returns either transaction (if one is active) or regular db connection.
func (s *DBStore) DB() ConnOrTx {
	if s.tx != nil {
		return s.tx
	}
	return s.db
}

// BeginTx starts a new transaction.
func (s *DBStore) BeginTx(ctx context.Context) error {
	if s.tx != nil {
		return errors.New("there is already a running transaction")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("beginning transaction: %v", err)
	}

	s.tx = tx
	return nil
}

// Commit commits and closes a transaction.
func (s *DBStore) Commit() error {
	if s.tx == nil {
		return errors.New("there is no running transaction")
	}

	err := s.tx.Commit()
	if err != nil {
		return fmt.Errorf("commiting transaction: %v", err)
	}

	s.tx = nil
	return nil
}

// Rollback rolls back and closes a transaction.
func (s *DBStore) Rollback() error {
	if s.tx == nil {
		// We have to assume that caller of Rollback knows what he is doing
		// nil will usually be the case when rollback is called in defer after the tx is commited
		return nil
	}

	err := s.tx.Rollback()
	if err != nil {
		return fmt.Errorf("rolling back transaction: %v", err)
	}

	s.tx = nil
	return nil
}
