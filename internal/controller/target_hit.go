package controller

import (
	"fmt"
	"strconv"
	"time"
)

func (i implementation) TargetWasHit(message string) (string, error) {
	i.logger.Debug("TargetWasHit controller, message: " + message)
	_, eventID, competitorID, extraParams, err := i.parseEvent(message)
	if err != nil {
		return "", fmt.Errorf("event 6: %w", err)
	}
	if eventID != 6 {
		return "", fmt.Errorf("event 6: expected EventID 6, got %d", eventID)
	}
	if len(extraParams) < 1 {
		return "", fmt.Errorf("event 6: missing target")
	}
	target, err := strconv.Atoi(extraParams[0])
	if err != nil {
		return "", fmt.Errorf("event 6: invalid target: %w", err)
	}
	if target <= 0 {
		return "", fmt.Errorf("event 6: target must be positive")
	}
	err = i.firingRangeUseCase.ShotTaken(competitorID, true)
	if err != nil {
		return "", fmt.Errorf("event 6: %w", err)
	}

	logMessage := fmt.Sprintf("[%s] The target(%d) has been hit by competitor(%d)",
		time.Now().Format("15:04:05.000"),
		target,
		competitorID)

	return logMessage, nil
}
