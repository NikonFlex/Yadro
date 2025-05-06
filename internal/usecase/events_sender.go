package usecase

import (
	"Solution/config"
	"fmt"
	"go.uber.org/zap"
	"os"
	"time"
)

type OutgoingEventsSender interface {
	SendEvent(event int) error
}

type sender struct {
	logger *zap.Logger
	config *config.Config
}

var _ OutgoingEventsSender = (*sender)(nil)

func NewSender(
	logger *zap.Logger,
	config *config.Config,
) *sender {
	return &sender{
		logger: logger,
		config: config,
	}
}

func (s sender) SendEvent(event int) error {
	s.logger.Debug("Sending event", zap.Int("event", event))
	newMsg := fmt.Sprintf("[%s] %d",
		time.Now().Format("15:04:05.000"),
		event)

	f, err := os.OpenFile(s.config.OutgoingEventsFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		s.logger.Error("Failed to open log file", zap.Error(err))
	}

	if _, err := f.WriteString(newMsg + "\n"); err != nil {
		s.logger.Error("Failed to write event message", zap.Error(err))
	}
	err = f.Close()
	if err != nil {
		s.logger.Error("Failed to close event file", zap.Error(err))
	}

	return nil
}
