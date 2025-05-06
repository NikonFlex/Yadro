package usecase

import (
	"Solution/config"
	"Solution/internal/entity"
	"fmt"
	"go.uber.org/zap"
	"os"
	"sort"
	"time"
)

func formatDuration(d time.Duration) string {
	minutes := int(d.Minutes())
	seconds := int(d.Seconds()) % 60
	milliseconds := (int(d.Milliseconds()) % 1000) / 10
	return fmt.Sprintf("%02d:%02d.%02d", minutes, seconds, milliseconds)
}

func calculateTotalTime(c *entity.Competitor) (string, time.Duration) {
	if c.Start.IsZero() {
		return "NotStarted", 0
	}
	if c.Finish.IsZero() {
		return "NotFinished", time.Hour * 24
	}

	total := c.Finish.Sub(c.Start)
	if !c.ScheduledStart.IsZero() {
		total += c.Start.Sub(c.ScheduledStart)
	}
	return formatDuration(total), total
}

func calculateLapInfo(lap *entity.Lap, distance float64) (string, string) {
	if lap.StartTime.IsZero() || lap.EndTime.IsZero() {
		return "-", "-"
	}
	duration := lap.EndTime.Sub(lap.StartTime)
	timeStr := formatDuration(duration)
	speed := distance / duration.Seconds()
	speedStr := fmt.Sprintf("%.3f", speed)
	return timeStr, speedStr
}

func generateCompetitorReport(c *entity.Competitor, cfg *config.Config) string {
	totalTimeStr, _ := calculateTotalTime(c)

	lapsStr := ""
	for _, lap := range c.Laps {
		timeStr, speedStr := calculateLapInfo(&lap, cfg.LapLen)
		lapsStr += fmt.Sprintf("{%s, %s} ", timeStr, speedStr)
	}
	if lapsStr == "" {
		lapsStr = "{} "
	} else {
		lapsStr = lapsStr[:len(lapsStr)-1]
	}

	penaltyLapsStr := ""
	for _, lap := range c.PenaltyLaps {
		timeStr, speedStr := calculateLapInfo(&lap, cfg.PenaltyLen)
		penaltyLapsStr += fmt.Sprintf("{%s, %s} ", timeStr, speedStr)
	}
	if penaltyLapsStr == "" {
		penaltyLapsStr = "{}"
	} else {
		penaltyLapsStr = penaltyLapsStr[:len(penaltyLapsStr)-1]
	}

	hits := 0
	shots := 0
	for _, fr := range c.FiringResults {
		hits += fr.Hits
		shots += cfg.Targets
	}
	firingStr := fmt.Sprintf("%d/%d", hits, shots)

	return fmt.Sprintf("[%s] %d %s Penalty:%s %s", totalTimeStr, c.ID, lapsStr, penaltyLapsStr, firingStr)
}

func generateReport(cfg *config.Config, logger *zap.Logger, competitors []entity.Competitor) error {
	sortedCompetitors := make([]entity.Competitor, len(competitors))
	copy(sortedCompetitors, competitors)

	sort.Slice(sortedCompetitors, func(i, j int) bool {
		_, timeI := calculateTotalTime(&sortedCompetitors[i])
		_, timeJ := calculateTotalTime(&sortedCompetitors[j])
		return timeI < timeJ
	})

	var reportLines []string
	for _, c := range sortedCompetitors {
		reportLines = append(reportLines, generateCompetitorReport(&c, cfg))
	}

	f, err := os.Create(cfg.ReportFileName)
	if err != nil {
		return fmt.Errorf("failed to create report file: %w", err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logger.Error("Failed to close report file", zap.Error(err))
		}
	}(f)

	for _, line := range reportLines {
		if _, err := f.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write to report file: %w", err)
		}
	}

	return nil
}
