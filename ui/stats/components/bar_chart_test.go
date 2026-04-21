package components

import (
	"strings"
	"testing"
	"time"

	"github.com/Bahaaio/pomo/db"
)

func TestBuildBars_RendersNonZeroDurationAsVisibleBar(t *testing.T) {
	chart := NewBarChart(12)
	stats := []db.DailyStat{
		{Date: "2026-02-02", WorkDuration: 8 * time.Hour, ScreenWorkDuration: 8 * time.Hour},
		{Date: "2026-02-08", WorkDuration: 45 * time.Minute, ScreenWorkDuration: 45 * time.Minute},
	}

	got := chart.View(stats)

	if !strings.Contains(got, barChar) {
		t.Fatalf("expected visible bar for non-zero duration, got %q", got)
	}

	if !strings.Contains(got, "45m") {
		t.Fatalf("expected duration label for bar, got %q", got)
	}
}

func TestSplitBarHeights_BothSourcesVisible(t *testing.T) {
	screen, other := splitBarHeights(50*time.Minute, 10*time.Minute, 6)
	if screen <= 0 || other <= 0 {
		t.Fatalf("expected both parts > 0, got screen=%d other=%d", screen, other)
	}
	if screen+other != 6 {
		t.Fatalf("expected total height 6, got %d", screen+other)
	}
}
