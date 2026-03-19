package gateway

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type StateStore interface {
	Load() (map[string]routeState, error)
	Save(key string, state routeState) error
	Close() error
}

type sqliteStateStore struct {
	db *sql.DB
}

func newSQLiteStateStore(dbPath string) StateStore {
	dbPath = strings.TrimSpace(dbPath)
	if dbPath == "" {
		return nil
	}
	db, err := sql.Open("sqlite", dbPath+"?_journal_mode=WAL")
	if err != nil {
		return nil
	}
	store := &sqliteStateStore{db: db}
	if err := store.init(); err != nil {
		_ = db.Close()
		return nil
	}
	return store
}

func (s *sqliteStateStore) init() error {
	if s == nil || s.db == nil {
		return fmt.Errorf("gateway state store is not configured")
	}
	_, err := s.db.Exec(`
	CREATE TABLE IF NOT EXISTS gateway_route_states (
		route_key TEXT PRIMARY KEY,
		day TEXT NOT NULL DEFAULT '',
		requests INTEGER NOT NULL DEFAULT 0,
		failures INTEGER NOT NULL DEFAULT 0,
		consecutive_failures INTEGER NOT NULL DEFAULT 0,
		breaker_state TEXT NOT NULL DEFAULT 'closed',
		breaker_open_until DATETIME,
		last_error TEXT NOT NULL DEFAULT '',
		last_decision_model TEXT NOT NULL DEFAULT '',
		updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		return fmt.Errorf("init gateway state store: %w", err)
	}
	if _, err := s.db.Exec(`ALTER TABLE gateway_route_states ADD COLUMN day TEXT NOT NULL DEFAULT ''`); err != nil && !strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
		return fmt.Errorf("migrate gateway state store: %w", err)
	}
	return nil
}

func (s *sqliteStateStore) Load() (map[string]routeState, error) {
	if s == nil || s.db == nil {
		return map[string]routeState{}, nil
	}
	rows, err := s.db.Query(`
		SELECT route_key, day, requests, failures, consecutive_failures, breaker_state, breaker_open_until, last_error, last_decision_model
		FROM gateway_route_states
	`)
	if err != nil {
		return nil, fmt.Errorf("load gateway state: %w", err)
	}
	defer rows.Close()

	states := make(map[string]routeState)
	for rows.Next() {
		var (
			key              string
			day              string
			requests         int
			failures         int
			consecutiveFails int
			breakerState     string
			breakerOpenUntil sql.NullTime
			lastError        string
			lastDecision     string
		)
		if err := rows.Scan(&key, &day, &requests, &failures, &consecutiveFails, &breakerState, &breakerOpenUntil, &lastError, &lastDecision); err != nil {
			return nil, fmt.Errorf("scan gateway state: %w", err)
		}
		state := routeState{
			Day:               day,
			Requests:          requests,
			Failures:          failures,
			ConsecutiveFails:  consecutiveFails,
			BreakerState:      breakerState,
			LastError:         lastError,
			LastDecisionModel: lastDecision,
		}
		if breakerOpenUntil.Valid {
			state.BreakerOpenUntil = breakerOpenUntil.Time.UTC()
		}
		states[key] = state
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate gateway state: %w", err)
	}
	return states, nil
}

func (s *sqliteStateStore) Save(key string, state routeState) error {
	if s == nil || s.db == nil || strings.TrimSpace(key) == "" {
		return nil
	}
	if state.Day == "" {
		state.Day = time.Now().UTC().Format("2006-01-02")
	}
	var breakerOpenUntil any
	if !state.BreakerOpenUntil.IsZero() {
		breakerOpenUntil = state.BreakerOpenUntil.UTC()
	}
	_, err := s.db.Exec(`
		INSERT INTO gateway_route_states (
			route_key, day, requests, failures, consecutive_failures, breaker_state, breaker_open_until, last_error, last_decision_model, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(route_key) DO UPDATE SET
			day = excluded.day,
			requests = excluded.requests,
			failures = excluded.failures,
			consecutive_failures = excluded.consecutive_failures,
			breaker_state = excluded.breaker_state,
			breaker_open_until = excluded.breaker_open_until,
			last_error = excluded.last_error,
			last_decision_model = excluded.last_decision_model,
			updated_at = excluded.updated_at
	`,
		key,
		state.Day,
		state.Requests,
		state.Failures,
		state.ConsecutiveFails,
		state.BreakerState,
		breakerOpenUntil,
		state.LastError,
		state.LastDecisionModel,
		time.Now().UTC(),
	)
	if err != nil {
		return fmt.Errorf("save gateway state %q: %w", key, err)
	}
	return nil
}

func (s *sqliteStateStore) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}
