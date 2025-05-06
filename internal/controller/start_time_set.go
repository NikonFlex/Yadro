package controller

import (
	"fmt"
	"time"
)

func (i implementation) StartTimeWasSet(message string) (string, error) {
	_, eventID, competitorID, extraParams, err := i.parseEvent(message)
	if err != nil {
		return "", fmt.Errorf("event 2: %w", err)
	}
	if eventID != 2 {
		return "", fmt.Errorf("event 2: expected EventID 2, got %d", eventID)
	}
	if len(extraParams) < 1 {
		return "", fmt.Errorf("event 2: missing start time")
	}
	startTime, err := time.Parse("15:04:05.000", extraParams[0])
	if err != nil {
		return "", fmt.Errorf("event 2: invalid start time format: %w", err)
	}
	err = i.competitorUseCase.SetStartTime(competitorID, startTime)
	if err != nil {
		return "", fmt.Errorf("event 2: %w", err)
	}

	logMessage := fmt.Sprintf("[%s] The start time for the competitor(%d) was set by a draw to %s",
		time.Now().Format("15:04:05.000"),
		competitorID,
		startTime.Format("15:04:05.000"))

	return logMessage, nil
}
