package controller

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

func (i implementation) parseEvent(line string) (time.Time, int, int, []string, error) {
	if !strings.HasPrefix(line, "[") {
		return time.Time{}, 0, 0, nil, fmt.Errorf("invalid event format: missing time prefix")
	}

	endTimeIdx := strings.Index(line, "]")
	if endTimeIdx == -1 {
		return time.Time{}, 0, 0, nil, fmt.Errorf("invalid event format: missing closing bracket")
	}

	timeStr := line[1:endTimeIdx]
	eventTime, err := time.Parse("15:04:05.000", timeStr)
	if err != nil {
		return time.Time{}, 0, 0, nil, fmt.Errorf("invalid time format: %w", err)
	}

	rest := strings.TrimSpace(line[endTimeIdx+1:])
	parts := strings.Fields(rest)
	if len(parts) < 2 {
		return time.Time{}, 0, 0, nil, fmt.Errorf("invalid event format: expected at least EventID and CompetitorID")
	}

	eventID, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, 0, 0, nil, fmt.Errorf("invalid EventID: %w", err)
	}

	competitorID, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, 0, 0, nil, fmt.Errorf("invalid CompetitorID: %w", err)
	}

	extraParams := parts[2:]
	return eventTime, eventID, competitorID, extraParams, nil
}
