package usecase

import (
	"Solution/config"
	"Solution/internal/entity"
	"Solution/internal/repo"
	"go.uber.org/zap"
	"time"
)

type LapUseCase interface {
	StartLap(competitorId int, startTime time.Time) error
	FinishLap(competitorId int, finishTime time.Time) error
	StartPenaltyLap(competitorId int, startTime time.Time) error
	FinishPenaltyLap(competitorId int, finishTime time.Time) error
}

type FiringRangeUseCase interface {
	StartFiringRange(competitorId int) error
	ShotTaken(competitorId int, result bool) error
	LeftFiringRange(competitorId int) error
}

type CompetitorUseCase interface {
	WriteReport() error
	Register(id int) error
	SetStartTime(id int, startTime time.Time) error
	Start(id int, startTime time.Time) error
	Finish(id int, finishTime time.Time) error
	CantContinue(id int, comment string) error
}

var _ LapUseCase = (*competition)(nil)
var _ FiringRangeUseCase = (*competition)(nil)
var _ CompetitorUseCase = (*competition)(nil)

type competition struct {
	logger                *zap.Logger
	config                *config.Config
	sender                OutgoingEventsSender
	competitorsRepository repo.CompetitorRepository
	lapRepository         repo.LapRepository
	firingRangeRepository repo.FiringRangeRepository
	stateMachine          *entity.StateMachine
}

func NewCompetition(
	logger *zap.Logger,
	config *config.Config,
	sender OutgoingEventsSender,
	competitorsRepository repo.CompetitorRepository,
	lapRepository repo.LapRepository,
	firingRangeRepository repo.FiringRangeRepository,
) *competition {
	return &competition{
		logger:                logger,
		config:                config,
		sender:                sender,
		competitorsRepository: competitorsRepository,
		lapRepository:         lapRepository,
		firingRangeRepository: firingRangeRepository,
		stateMachine:          entity.NewStateMachine(),
	}
}

func (c competition) moveToState(competitorId int, next entity.Status) error {
	status, err := c.competitorsRepository.GetStatus(competitorId)
	if err != nil {
		return err
	}

	if ok := c.stateMachine.MoveTo(status, next); !ok {
		c.logger.Error("Transition not accepted", zap.Int("current", int(status)), zap.Int("next", int(next)))
		return entity.ErrStateIsNotAccepted
	}

	err = c.competitorsRepository.SetStatus(competitorId, next)
	if err != nil {
		return err
	}

	return nil
}
