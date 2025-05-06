package usecase

import (
	"Solution/internal/entity"
	"go.uber.org/zap"
)

func (c competition) StartFiringRange(competitorId int) error {
	c.logger.Debug("Competitor start firing range", zap.Int("competitorId", competitorId))
	err := c.moveToState(competitorId, entity.StatusOnFiringRange)
	if err != nil {
		return err
	}

	err = c.firingRangeRepository.StartFiringRange(competitorId)
	if err != nil {
		return err
	}

	return nil
}

func (c competition) LeftFiringRange(competitorId int) error {
	c.logger.Debug("Competitor start firing range", zap.Int("competitorId", competitorId))
	competitor, err := c.competitorsRepository.Get(competitorId)
	if err != nil {
		return err
	}

	lastFiringResult := competitor.FiringResults[len(competitor.FiringResults)-1]
	if lastFiringResult.Hits != c.config.Targets {
		c.logger.Debug("Competitor missed some shots", zap.Int("competitorId", competitorId))
		err := c.moveToState(competitorId, entity.StatusOnPenaltyLap)
		if err != nil {
			return err
		}
	} else {
		c.logger.Debug("Competitor hits all targets", zap.Int("competitorId", competitorId))
		err := c.moveToState(competitorId, entity.StatusOnLap)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c competition) ShotTaken(competitorId int, result bool) error {
	c.logger.Debug("Competitor shot taken", zap.Int("id", competitorId), zap.Bool("result", result))
	status, err := c.competitorsRepository.GetStatus(competitorId)
	if status != entity.StatusOnFiringRange {
		c.logger.Error("Competitor wasn't on firing range", zap.Int("competitorId", competitorId))
		return entity.ErrStateIsNotAccepted
	}

	err = c.firingRangeRepository.ShotTaken(competitorId, result)
	if err != nil {
		return err
	}

	return nil
}
