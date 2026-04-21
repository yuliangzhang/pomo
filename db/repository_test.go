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

func TestGetDailyStats_SplitsScreenAndOther(t *testing.T) {
	repo := newTestRepo(t)
	day := time.Date(2026, 2, 14, 9, 0, 0, 0, time.Local)

	if err := repo.CreateSessionWithSource(day, time.Hour, WorkSession, ScreenSource); err != nil {
		t.Fatalf("create screen session: %v", err)
	}
	if err := repo.CreateSessionWithSource(day, 27*time.Minute, WorkSession, OtherSource); err != nil {
		t.Fatalf("create other session: %v", err)
	}

	stats, err := repo.getDailyStats(day, day)
	if err != nil {
		t.Fatalf("get daily stats: %v", err)
	}
	if len(stats) != 1 {
		t.Fatalf("expected 1 stat row, got %d", len(stats))
	}

	if stats[0].ScreenWorkDuration != time.Hour {
		t.Fatalf("screen duration = %v, want %v", stats[0].ScreenWorkDuration, time.Hour)
	}
	if stats[0].OtherWorkDuration != 27*time.Minute {
		t.Fatalf("other duration = %v, want %v", stats[0].OtherWorkDuration, 27*time.Minute)
	}
	if stats[0].WorkDuration != (time.Hour + 27*time.Minute) {
		t.Fatalf("work duration = %v, want %v", stats[0].WorkDuration, time.Hour+27*time.Minute)
	}
}
