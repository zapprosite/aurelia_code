package agent

import (
	"database/sql"
	"fmt"
	"time"

	_ "modernc.org/sqlite"
)

type SQLiteTaskStore struct {
	db            *sql.DB
	leaseDuration time.Duration
}

func NewSQLiteTaskStore(dbPath string) (*SQLiteTaskStore, error) {
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("open sqlite task store: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)

	store := &SQLiteTaskStore{
		db:            db,
		leaseDuration: 30 * time.Second,
	}
	if err := store.initialize(); err != nil {
		_ = db.Close()
		return nil, err
	}

	return store, nil
}

// NewTeamManager returns a TeamManager interface powered by the SQLiteTaskStore.
func NewTeamManager(s *SQLiteTaskStore) (TeamManager, error) {
	if s == nil {
		return nil, fmt.Errorf("task store is nil")
	}
	return s, nil
}

func (s *SQLiteTaskStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
