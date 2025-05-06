package controller

import (
	"fmt"
	"time"
)

func (i implementation) EndedMainLap(message string) (string, error) {
	i.logger.Debug("EndMainLoop controller, message: " + message)
	eventTime, eventID, competitorID, extraParams, err := i.parseEvent(message)
	if err != nil {
		return "", fmt.Errorf("event 10: %w", err)
	}
	if eventID != 10 {
		return "", fmt.Errorf("event 10: expected EventID 10, got %d", eventID)
	}
	if len(extraParams) > 0 {
		return "", fmt.Errorf("event 10: unexpected extra parameters: %v", extraParams)
	}

	err = i.lapUseCase.FinishLap(competitorID, eventTime)
	if err != nil {
		return "", fmt.Errorf("event 10: %w", err)
	}

	logMessage := fmt.Sprintf("[%s] The competitor(%d) ended the main lap",
		time.Now().Format("15:04:05.000"),
		competitorID)

	return logMessage, nil
}
