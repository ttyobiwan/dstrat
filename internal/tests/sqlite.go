package tests

import (
	"database/sql"
	"os"
	"testing"

	"github.com/ttyobiwan/dstrat/internal/sqlite"
)

func GetTestDB(t *testing.T) *sql.DB {
	db, err := sqlite.GetDB("default_test.sqlite")
	if err != nil {
		t.Fatalf("getting db: %v", err)
	}
	err = sqlite.Configure(db)
	if err != nil {
		t.Fatalf("configuring db: %v", err)
	}
	err = sqlite.Migrate(db)
	if err != nil {
		t.Fatalf("migrating db: %v", err)
	}
	t.Cleanup(func() {
		if err := db.Close(); err != nil {
			t.Logf("closing db: %v", err)
		}
		if err := os.Remove("default_test.sqlite"); err != nil {
			t.Errorf("removing file: %v", err)
		}
	})
	return db
}
