package controller

import (
	"fmt"
	"strings"
	"time"
)

func (i implementation) CantContinue(message string) (string, error) {
	i.logger.Debug("CantContinue controller, message: " + message)
	_, eventID, competitorID, extraParams, err := i.parseEvent(message)
	if err != nil {
		return "", fmt.Errorf("event 11: %w", err)
	}
	if eventID != 11 {
		return "", fmt.Errorf("event 11: expected EventID 11, got %d", eventID)
	}
	if len(extraParams) < 1 {
		return "", fmt.Errorf("event 11: missing comment")
	}
	comment := strings.Join(extraParams, " ")

	err = i.competitorUseCase.CantContinue(competitorID, comment)
	if err != nil {
		return "", fmt.Errorf("event 11: %w", err)
	}

	logMessage := fmt.Sprintf("[%s] The competitor(%d) can't continue: %s",
		time.Now().Format("15:04:05.000"),
		competitorID,
		comment)

	return logMessage, nil
}
