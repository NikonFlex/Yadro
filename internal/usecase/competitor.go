package usecase

import (
	"Solution/internal/entity"
	"go.uber.org/zap"
	"time"
)

func (c competition) Register(id int) error {
	c.logger.Debug("Registering competitor", zap.Int("id", id))
	err := c.competitorsRepository.Register(id)
	if err != nil {
		return err
	}

	err = c.competitorsRepository.SetStatus(id, entity.StatusRegistered)
	if err != nil {
		return err
	}

	return nil
}

func (c competition) WriteReport() error {
	err := generateReport(c.config, c.logger, c.competitorsRepository.GetAll())
	if err != nil {
		c.logger.Error("Failed to generate report", zap.Error(err))
		return err
	}

	return nil
}

func (c competition) SetStartTime(id int, startTime time.Time) error {
	c.logger.Debug("Setting start time competitor", zap.Int("id", id))
	status, err := c.competitorsRepository.GetStatus(id)
	if err != nil {
		return err
	}
	if status != entity.StatusRegistered {
		c.logger.Error("Competitor is not registered", zap.Int("id", id))
		return entity.ErrStateIsNotAccepted
	}
	err = c.competitorsRepository.SetStartTime(id, startTime)
	if err != nil {
		return err
	}

	return nil
}

func (c competition) Start(id int, startTime time.Time) error {
	c.logger.Debug("Starting competitor", zap.Int("id", id))
	competitor, err := c.competitorsRepository.Get(id)
	if err != nil {
		return err
	}

	if !c.startedInInterval(competitor.ScheduledStart, competitor.ScheduledStart.Add(c.config.StartDelta), startTime) {
		c.logger.Info("Competitor is not started in interval ", zap.Time("Scheduled", competitor.ScheduledStart), zap.Time("Actual", startTime), zap.Int("id", id))
		err := c.moveToState(id, entity.StatusDisqualified)
		if err != nil {
			return err
		}

		err = c.sender.SendEvent(32)
		if err != nil {
			return err
		}

		return nil
	}

	err = c.moveToState(id, entity.StatusOnLap)
	if err != nil {
		return err
	}

	err = c.competitorsRepository.Start(id, startTime)
	if err != nil {
		return err
	}

	err = c.lapRepository.StartLap(id, startTime)
	if err != nil {
		return err
	}

	return nil
}

func (c competition) Finish(id int, finishTime time.Time) error {
	c.logger.Debug("Finishing competitor", zap.Int("id", id))
	err := c.moveToState(id, entity.StatusFinished)
	if err != nil {
		return err
	}

	err = c.competitorsRepository.Finish(id, finishTime)
	if err != nil {
		return err
	}

	err = c.sender.SendEvent(33)
	if err != nil {
		return err
	}

	return nil
}

func (c competition) CantContinue(id int, comment string) error {
	c.logger.Debug("Competitor cant continue", zap.Int("id", id))
	err := c.moveToState(id, entity.StatusNotFinished)
	if err != nil {
		return err
	}

	err = c.competitorsRepository.CantContinue(id, comment)
	if err != nil {
		return err
	}

	return nil
}

func (c competition) startedInInterval(start time.Time, end time.Time, value time.Time) bool {
	return start.Before(value) && end.After(value)
}
