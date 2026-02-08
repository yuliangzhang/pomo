package stats

import (
	"testing"
	"time"

	"github.com/Bahaaio/pomo/db"
)

func TestBuildTodayWorkLine(t *testing.T) {
	now := time.Date(2026, 2, 8, 10, 0, 0, 0, time.Local)

	stats := []db.DailyStat{
		{Date: "2026-02-07", WorkDuration: 40 * time.Minute},
		{Date: "2026-02-08", WorkDuration: 95 * time.Minute},
	}

	got := buildTodayWorkLine(stats, now)
	want := "today work 1h35m"

	if got != want {
		t.Fatalf("buildTodayWorkLine() = %q, want %q", got, want)
	}
}

func TestBuildTodayWorkLine_NoTodayData(t *testing.T) {
	now := time.Date(2026, 2, 8, 10, 0, 0, 0, time.Local)

	stats := []db.DailyStat{
		{Date: "2026-02-07", WorkDuration: 40 * time.Minute},
	}

	got := buildTodayWorkLine(stats, now)
	want := "today work 0m"

	if got != want {
		t.Fatalf("buildTodayWorkLine() = %q, want %q", got, want)
	}
}

