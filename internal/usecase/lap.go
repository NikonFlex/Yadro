package usecase

import (
	"Solution/internal/entity"
	"go.uber.org/zap"
	"time"
)

func (c competition) StartLap(competitorId int, startTime time.Time) error {
	c.logger.Debug("Competitor starts lap", zap.Int("id", competitorId))
	err := c.moveToState(competitorId, entity.StatusOnLap)
	if err != nil {
		return err
	}

	err = c.lapRepository.StartLap(competitorId, startTime)
	if err != nil {
		return err
	}

	return nil
}

func (c competition) FinishLap(competitorId int, finishTime time.Time) error {
	c.logger.Debug("Competitor finish lap", zap.Int("id", competitorId))
	competitor, err := c.competitorsRepository.Get(competitorId)
	if err != nil {
		return err
	}

	err = c.lapRepository.FinishLap(competitorId, finishTime)
	if err != nil {
		return err
	}

	if c.config.Laps == len(competitor.Laps) {
		c.logger.Debug("All laps finished", zap.Int("id", competitorId))
		err := c.moveToState(competitorId, entity.StatusFinished)
		if err != nil {
			return err
		}

		err = c.competitorsRepository.Finish(competitorId, finishTime)
		if err != nil {
			return err
		}
	} else {
		c.logger.Debug("Go on next lap", zap.Int("id", competitorId))
		err := c.StartLap(competitorId, finishTime)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c competition) StartPenaltyLap(competitorId int, startTime time.Time) error {
	c.logger.Debug("Competitor starts penalty lap", zap.Int("id", competitorId))
	err := c.moveToState(competitorId, entity.StatusOnPenaltyLap)
	if err != nil {
		return err
	}

	err = c.lapRepository.StartPenaltyLap(competitorId, startTime)
	if err != nil {
		return err
	}

	return nil
}

func (c competition) FinishPenaltyLap(competitorId int, finishTime time.Time) error {
	c.logger.Debug("Competitor finished penalty lap", zap.Int("id", competitorId))
	err := c.moveToState(competitorId, entity.StatusOnLap)
	if err != nil {
		return err
	}

	err = c.lapRepository.FinishPenaltyLap(competitorId, finishTime)
	if err != nil {
		return err
	}

	return nil
}
