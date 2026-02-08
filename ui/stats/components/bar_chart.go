// Package components implements UI components for stats.
package components

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bahaaio/pomo/db"
	"github.com/Bahaaio/pomo/ui/colors"
	"github.com/charmbracelet/lipgloss"
)

var barStyle = lipgloss.NewStyle().Foreground(colors.WorkSessionFg)

const (
	barChar     = "█"
	axisChar    = "│"
	tickChar    = "┤"
	cornerChar  = "└"
	lineChar    = "─"
	paddingChar = " "

	minBarWidth = 3
	spacing     = 2
	daysInWeek  = 7
)

var spacer = strings.Repeat(paddingChar, spacing)

type chartLayout struct {
	barHeight       int
	barWidth        int
	yAxisLabelWidth int
	yAxisWidth      int // label + space + tick char
	barAreaWidth    int
	totalWidth      int
}

type BarChart struct {
	chartLayout
}

func NewBarChart(height int) BarChart {
	return BarChart{
		chartLayout: chartLayout{
			barHeight: height - 1 - 1, // leave space for x-axis and labels
		},
	}
}

func (b *BarChart) calculateLayout(stats []db.DailyStat, maxDuration, scale time.Duration) chartLayout {
	longestLabel := 0

	for duration := maxDuration; duration > 0; duration -= scale {
		label := formatDuration(duration)
		longestLabel = max(longestLabel, len(label))
	}

	barWidth := getBarWidth(stats)
	yAxisLabelWidth := longestLabel
	yAxisWidth := yAxisLabelWidth + 1 + 1 // length of label + space + tick char

	barAreaWidth := spacing + (barWidth+spacing)*daysInWeek

	return chartLayout{
		barHeight:       b.barHeight,
		barWidth:        barWidth,
		yAxisLabelWidth: yAxisLabelWidth,
		yAxisWidth:      yAxisWidth,
		barAreaWidth:    barAreaWidth,
		totalWidth:      yAxisWidth + barAreaWidth,
	}
}

func (b *BarChart) View(stats []db.DailyStat) string {
	if len(stats) == 0 {
		return ""
	}

	maxDuration := getMaxDuration(stats)

	// dividing by half of max height to leave space for tick chars
	targetTicks := b.barHeight / 2
	scale := calculateScale(maxDuration, targetTicks)

	// fallback for empty stats
	if maxDuration == 0 {
		maxDuration = time.Hour
		scale = time.Minute * 10
	}

	b.chartLayout = b.calculateLayout(stats, maxDuration, scale)

	yAxis := b.buildYAxis(maxDuration, scale)
	bars := b.buildBars(stats, maxDuration)

	top := lipgloss.JoinHorizontal(lipgloss.Left, yAxis, spacer, bars)
	xAxis := b.buildXAxis()
	labels := b.buildLabels(stats)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		top,
		xAxis,
		labels,
	)
}

func (b *BarChart) buildBars(stats []db.DailyStat, maxDuration time.Duration) string {
	bars := make([]string, 0, len(stats))
	totalRows := b.barHeight + 1 // extra row for value labels

	for _, stat := range stats {
		barHeight := 0
		if stat.WorkDuration > 0 {
			barHeight = int((float64(stat.WorkDuration) / float64(maxDuration)) * float64(b.barHeight))
			// Ensure any non-zero duration is visible in the chart.
			if barHeight == 0 {
				barHeight = 1
			}
		}

		bar := renderBar(totalRows, barHeight, b.barWidth, formatDurationLabel(stat.WorkDuration))
		bars = append(bars, bar, spacer)
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, bars...)
}

func (b *BarChart) buildYAxis(maxDuration, scale time.Duration) string {
	builder := strings.Builder{}

	// smallest duration to print
	epsilon := time.Millisecond * 500

	// print all rows with same width of longest duration label
	for duration := maxDuration; duration >= epsilon; duration -= scale {
		tick := fmt.Sprintf("%-*s %s\n", b.yAxisLabelWidth, formatDuration(duration), tickChar)
		builder.WriteString(tick)

		axis := strings.Repeat(paddingChar, b.yAxisLabelWidth) + paddingChar + axisChar
		builder.WriteString(axis)

		// don't print the last new line
		if duration-scale >= epsilon {
			builder.WriteString("\n")
		}
	}

	return builder.String()
}

func (b *BarChart) buildXAxis() string {
	zeroLabel := fmt.Sprintf("%-*s", b.yAxisLabelWidth, "0")
	return zeroLabel + " " + cornerChar + strings.Repeat(lineChar, b.barAreaWidth)
}

func (b *BarChart) buildLabels(stats []db.DailyStat) string {
	var labels strings.Builder

	for _, stat := range stats {
		day := getDayLabel(stat.Date, b.barWidth)
		labels.WriteString(day)
		labels.WriteString(spacer)
	}

	// yaxis width + spacing between yaxis and bars
	paddingLength := b.yAxisWidth + spacing
	padding := strings.Repeat(paddingChar, paddingLength)

	return padding + labels.String()
}

func getDayLabel(day string, width int) string {
	t, err := time.Parse(db.DateFormat, day)
	if err != nil {
		return strings.Repeat(paddingChar, width)
	}

	// get first three letters of weekday
	return centerText(t.Weekday().String()[:3], width)
}

func renderBar(totalRows, barHeight, barWidth int, label string) string {
	if totalRows <= 0 {
		return ""
	}

	blank := strings.Repeat(paddingChar, barWidth)
	rows := make([]string, totalRows)
	for i := range rows {
		rows[i] = blank
	}

	if barHeight <= 0 {
		// no work duration: render label on the last row (just above x-axis)
		rows[totalRows-1] = centerText(label, barWidth)
		return strings.Join(rows, "\n")
	}

	barHeight = min(barHeight, totalRows-1)
	labelRow := totalRows - barHeight - 1
	rows[labelRow] = centerText(label, barWidth)

	filledRow := barStyle.Render(strings.Repeat(barChar, barWidth))
	for i := labelRow + 1; i <= labelRow+barHeight && i < totalRows; i++ {
		rows[i] = filledRow
	}

	return strings.Join(rows, "\n")
}

func getBarWidth(stats []db.DailyStat) int {
	width := minBarWidth

	for _, stat := range stats {
		width = max(width, len(formatDurationLabel(stat.WorkDuration)))
	}

	return width
}

func formatDurationLabel(d time.Duration) string {
	if d <= 0 {
		return "0m"
	}

	if d < time.Minute {
		seconds := int(d.Seconds())
		if seconds == 0 {
			seconds = 1
		}
		return fmt.Sprintf("%ds", seconds)
	}

	minutes := int(d.Minutes())
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}

	hours := minutes / 60
	remainingMinutes := minutes % 60
	if remainingMinutes == 0 {
		return fmt.Sprintf("%dh", hours)
	}

	return fmt.Sprintf("%dh%dm", hours, remainingMinutes)
}

func centerText(text string, width int) string {
	if width <= 0 {
		return ""
	}

	if len(text) >= width {
		return text[:width]
	}

	left := (width - len(text)) / 2
	right := width - len(text) - left

	return strings.Repeat(paddingChar, left) + text + strings.Repeat(paddingChar, right)
}

func getMaxDuration(stats []db.DailyStat) time.Duration {
	var maxDuration time.Duration

	for _, stat := range stats {
		if stat.WorkDuration > maxDuration {
			maxDuration = stat.WorkDuration
		}
	}

	return maxDuration
}

func calculateScale(maxDuration time.Duration, targetTicks int) time.Duration {
	scale := time.Duration(
		float64(maxDuration.Milliseconds())/float64(targetTicks),
	) * time.Millisecond

	// minimum scale of 100 ms to avoid too many ticks
	scale = max(scale, time.Millisecond*100)

	return scale
}

func formatDuration(d time.Duration) string {
	seconds := d.Seconds()
	if seconds < 60 {
		if seconds == float64(int(seconds)) {
			return fmt.Sprintf("%ds", int(seconds))
		}

		// show one decimal place for seconds less than 1
		return fmt.Sprintf("%0.1fs", seconds)
	}

	minutes := int(d.Minutes())
	if minutes < 60 {
		return fmt.Sprintf("%dm", minutes)
	}

	hours := minutes / 60
	mins := minutes % 60

	if mins == 0 {
		return fmt.Sprintf("%dh", hours)
	}

	return fmt.Sprintf("%dh%dm", hours, mins)
}
