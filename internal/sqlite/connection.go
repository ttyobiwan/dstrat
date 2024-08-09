package sqlite

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

var (
	initialMigration = `CREATE TABLE IF NOT EXISTS users (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	username VARCHAR(127) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS topics (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name VARCHAR(255) UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	title VARCHAR(255) NOT NULL,
	content TEXT NOT NULL,
	author_id INTEGER NOT NULL,
	FOREIGN KEY (author_id) REFERENCES users(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS post_topics (
	post_id INTEGER,
	topic_id INTEGER,
	PRIMARY KEY (post_id, topic_id),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS topic_followers (
	topic_id INTEGER,
	user_id INTEGER,
	PRIMARY KEY (topic_id, user_id),
	FOREIGN KEY (topic_id) REFERENCES topics(id) ON DELETE CASCADE,
	FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);`
)

// GetDB creates a new SQLite connection.
func GetDB(filename string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, fmt.Errorf("opening sqlite db: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("pinging sqlite db: %v", err)
	}

	return db, nil
}

// Configure adds additional configuration for the SQLite database and connection.
func Configure(db *sql.DB) error {
	_, err := db.Exec("PRAGMA journal_mode=WAL")
	if err != nil {
		return fmt.Errorf("enabling wal mode: %v", err)
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON")
	if err != nil {
		return fmt.Errorf("enabling fks: %v", err)
	}

	return nil
}

// Migrate runs the database migrations.
func Migrate(db *sql.DB) error {
	_, err := db.Exec(initialMigration)
	if err != nil {
		return fmt.Errorf("migrating sqlite db: %v", err)
	}
	return nil
}
