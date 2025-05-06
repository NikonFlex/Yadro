package controller

import (
	"fmt"
	"time"
)

func (i implementation) LeftFiringRange(message string) (string, error) {
	i.logger.Debug("LeftFiringRange controller, message: " + message)
	_, eventID, competitorID, extraParams, err := i.parseEvent(message)
	if err != nil {
		return "", fmt.Errorf("event 7: %w", err)
	}
	if eventID != 7 {
		return "", fmt.Errorf("event 7: expected EventID 7, got %d", eventID)
	}
	if len(extraParams) > 0 {
		return "", fmt.Errorf("event 7: unexpected extra parameters: %v", extraParams)
	}
	err = i.firingRangeUseCase.LeftFiringRange(competitorID)
	if err != nil {
		return "", fmt.Errorf("firing range 7: %w", err)
	}

	logMessage := fmt.Sprintf("[%s] The competitor(%d) left the firing range",
		time.Now().Format("15:04:05.000"),
		competitorID)

	return logMessage, nil
}
