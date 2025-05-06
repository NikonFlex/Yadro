package controller

import (
	"fmt"
	"time"
)

func (i implementation) OnPenaltyLap(message string) (string, error) {
	i.logger.Debug("OnPenaltyLap controller, message: " + message)
	eventTime, eventID, competitorID, extraParams, err := i.parseEvent(message)
	if err != nil {
		return "", fmt.Errorf("event 8: %w", err)
	}
	if eventID != 8 {
		return "", fmt.Errorf("event 8: expected EventID 8, got %d", eventID)
	}
	if len(extraParams) > 0 {
		return "", fmt.Errorf("event 8: unexpected extra parameters: %v", extraParams)
	}
	err = i.lapUseCase.StartPenaltyLap(competitorID, eventTime)
	if err != nil {
		return "", fmt.Errorf("event 8: %w", err)
	}

	logMessage := fmt.Sprintf("[%s] The competitor(%d) entered the penalty laps",
		time.Now().Format("15:04:05.000"),
		competitorID)

	return logMessage, nil
}
