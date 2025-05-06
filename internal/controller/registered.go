package controller

import (
	"fmt"
	"time"
)

func (i implementation) Registered(message string) (string, error) {
	i.logger.Debug("Registered controller, message: " + message)
	_, eventID, competitorID, extraParams, err := i.parseEvent(message)
	if err != nil {
		return "", fmt.Errorf("event 1: %w", err)
	}
	if eventID != 1 {
		return "", fmt.Errorf("event 1: expected EventID 1, got %d", eventID)
	}
	if len(extraParams) > 0 {
		return "", fmt.Errorf("event 1: unexpected extra parameters: %v", extraParams)
	}
	err = i.competitorUseCase.Register(competitorID)
	if err != nil {
		return "", fmt.Errorf("event 2: %w", err)
	}

	logMessage := fmt.Sprintf("[%s] The competitor(%d) registered",
		time.Now().Format("15:04:05.000"),
		competitorID)

	return logMessage, nil
}
