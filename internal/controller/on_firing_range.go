package controller

import (
	"fmt"
	"strconv"
	"time"
)

func (i implementation) OnFiringRange(message string) (string, error) {
	i.logger.Debug("OnFiringRange controller, message: " + message)
	_, eventID, competitorID, extraParams, err := i.parseEvent(message)
	if err != nil {
		return "", fmt.Errorf("event 5: %w", err)
	}
	if eventID != 5 {
		return "", fmt.Errorf("event 5: expected EventID 5, got %d", eventID)
	}
	if len(extraParams) < 1 {
		return "", fmt.Errorf("event 5: missing firing range")
	}
	firingRange, err := strconv.Atoi(extraParams[0])
	if err != nil {
		return "", fmt.Errorf("event 5: invalid firing range: %w", err)
	}
	if firingRange <= 0 {
		return "", fmt.Errorf("event 5: firing range must be positive")
	}
	err = i.firingRangeUseCase.StartFiringRange(competitorID)
	if err != nil {
		return "", fmt.Errorf("event 5: %w", err)
	}

	logMessage := fmt.Sprintf("[%s] The competitor(%d) is on firing range(%d)",
		time.Now().Format("15:04:05.000"),
		competitorID,
		firingRange)

	return logMessage, nil
}
