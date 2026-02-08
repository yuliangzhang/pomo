package ui

import (
	"testing"
	"time"
)

func TestCalculateSessionStartTime(t *testing.T) {
	now := time.Date(2026, 2, 8, 0, 10, 0, 0, time.Local)
	elapsed := 25 * time.Minute

	got := calculateSessionStartTime(now, elapsed)
	want := time.Date(2026, 2, 7, 23, 45, 0, 0, time.Local)

	if !got.Equal(want) {
		t.Fatalf("calculateSessionStartTime() = %v, want %v", got, want)
	}
}

func TestCalculateSessionStartTime_NonPositiveElapsed(t *testing.T) {
	now := time.Date(2026, 2, 8, 0, 10, 0, 0, time.Local)

	got := calculateSessionStartTime(now, 0)
	if !got.Equal(now) {
		t.Fatalf("calculateSessionStartTime() with 0 elapsed = %v, want %v", got, now)
	}
}

