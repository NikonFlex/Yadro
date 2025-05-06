package controller

import (
	"fmt"
	"time"
)

func (i implementation) LeftPenaltyLap(message string) (string, error) {
	i.logger.Debug("LeftPenaltyLap controller, message: " + message)
	eventTime, eventID, competitorID, extraParams, err := i.parseEvent(message)
	if err != nil {
		return "", fmt.Errorf("event 9: %w", err)
	}
	if eventID != 9 {
		return "", fmt.Errorf("event 9: expected EventID 9, got %d", eventID)
	}
	if len(extraParams) > 0 {
		return "", fmt.Errorf("event 9: unexpected extra parameters: %v", extraParams)
	}

	err = i.lapUseCase.FinishPenaltyLap(competitorID, eventTime)
	logMessage := fmt.Sprintf("[%s] The competitor(%d) left the penalty laps",
		time.Now().Format("15:04:05.000"),
		competitorID)

	return logMessage, nil
}
