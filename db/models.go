package db

import (
	"time"

	"github.com/Bahaaio/pomo/config"
)

var schema = `
CREATE TABLE IF NOT EXISTS sessions(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	type TEXT NOT NULL,
	duration INTEGER NOT NULL,
	started_at TEXT NOT NULL,
	source TEXT NOT NULL DEFAULT 'screen'
);
`

type Session struct {
	ID        int           `db:"id"`
	Type      string        `db:"type"`
	Duration  time.Duration `db:"duration"`
	StartedAt time.Time     `db:"started_at"`
	Source    string        `db:"source"`
}

type AllTimeStats struct {
	TotalSessions      int           `db:"total_sessions"`
	TotalWorkDuration  time.Duration `db:"total_work_duration"`
	TotalBreakDuration time.Duration `db:"total_break_duration"`
}

type DailyStat struct {
	Date               string        `db:"day"`
	ScreenWorkDuration time.Duration `db:"screen_work_duration"`
	OtherWorkDuration  time.Duration `db:"other_work_duration"`
	WorkDuration       time.Duration `db:"work_duration"`
}

type StreakStats struct {
	Current int
	Best    int
}

type SessionType string
type SessionSource string

const (
	WorkSession  SessionType = "work"
	BreakSession SessionType = "break"

	ScreenSource SessionSource = "screen"
	OtherSource  SessionSource = "other"
)

func GetSessionType(taskType config.TaskType) SessionType {
	if taskType == config.WorkTask {
		return WorkSession
	}
	return BreakSession
}
