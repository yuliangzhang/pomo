package db

import (
	"database/sql"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
)

func newTestRepo(t *testing.T) *SessionRepo {
	t.Helper()

	database, err := sqlx.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if _, err := database.Exec(schema); err != nil {
		t.Fatalf("create schema: %v", err)
	}

	return NewSessionRepo(database)
}

func TestExtendLatestSession(t *testing.T) {
	repo := newTestRepo(t)
	start := time.Date(2026, 2, 13, 9, 0, 0, 0, time.Local)

	if err := repo.CreateSession(start, time.Hour, WorkSession); err != nil {
		t.Fatalf("create session: %v", err)
	}

	if err := repo.ExtendLatestSession(27*time.Minute, WorkSession); err != nil {
		t.Fatalf("extend latest session: %v", err)
	}

	stats, err := repo.GetAllTimeStats()
	if err != nil {
		t.Fatalf("get all-time stats: %v", err)
	}

	if stats.TotalWorkDuration != (time.Hour + 27*time.Minute) {
		t.Fatalf("total work duration = %v, want %v", stats.TotalWorkDuration, time.Hour+27*time.Minute)
	}
}

func TestExtendLatestSession_NoRows(t *testing.T) {
	repo := newTestRepo(t)

	err := repo.ExtendLatestSession(10*time.Minute, WorkSession)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}

	if err != sql.ErrNoRows {
		t.Fatalf("error = %v, want %v", err, sql.ErrNoRows)
	}
}
