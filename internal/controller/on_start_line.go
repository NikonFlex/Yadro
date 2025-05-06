package controller

import (
	"fmt"
	"time"
)

func (i implementation) OnStartLine(message string) (string, error) {
	i.logger.Debug("OnStartLine controller, message: " + message)
	// not used

	_, _, competitorID, _, err := i.parseEvent(message)
	if err != nil {
		return "", fmt.Errorf("event 8: %w", err)
	}

	logMessage := fmt.Sprintf("[%s] The competitor(%d) ended the main lap",
		time.Now().Format("15:04:05.000"),
		competitorID)

	return logMessage, nil
}
