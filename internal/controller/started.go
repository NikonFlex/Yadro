package controller

import (
	"fmt"
	"time"
)

func (i implementation) Started(message string) (string, error) {
	i.logger.Debug("Started controller, message: " + message)
	eventTime, eventID, competitorID, extraParams, err := i.parseEvent(message)
	if err != nil {
		return "", fmt.Errorf("event 4: %w", err)
	}
	if eventID != 4 {
		return "", fmt.Errorf("event 4: expected EventID 4, got %d", eventID)
	}
	if len(extraParams) > 0 {
		return "", fmt.Errorf("event 4: unexpected extra parameters: %v", extraParams)
	}
	err = i.competitorUseCase.Start(competitorID, eventTime)
	if err != nil {
		return "", fmt.Errorf("event 5: %w", err)
	}

	logMessage := fmt.Sprintf("[%s] The competitor(%d) has started",
		time.Now().Format("15:04:05.000"),
		competitorID)

	return logMessage, nil
}
