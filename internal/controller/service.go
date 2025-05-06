package controller

import (
	"Solution/internal/usecase"
	"go.uber.org/zap"
)

type Controller interface {
	Registered(message string) (string, error)
	StartTimeWasSet(message string) (string, error)
	OnStartLine(message string) (string, error)
	Started(message string) (string, error)
	OnFiringRange(message string) (string, error)
	TargetWasHit(message string) (string, error)
	LeftFiringRange(message string) (string, error)
	OnPenaltyLap(message string) (string, error)
	LeftPenaltyLap(message string) (string, error)
	EndedMainLap(message string) (string, error)
	CantContinue(message string) (string, error)
	WriteReport() error
}

var _ Controller = (*implementation)(nil)

type implementation struct {
	logger             *zap.Logger
	lapUseCase         usecase.LapUseCase
	competitorUseCase  usecase.CompetitorUseCase
	firingRangeUseCase usecase.FiringRangeUseCase
}

func New(
	logger *zap.Logger,
	lapUseCase usecase.LapUseCase,
	competitorUseCase usecase.CompetitorUseCase,
	firingRangeUseCase usecase.FiringRangeUseCase,
) *implementation {
	return &implementation{
		logger:             logger,
		lapUseCase:         lapUseCase,
		competitorUseCase:  competitorUseCase,
		firingRangeUseCase: firingRangeUseCase,
	}
}
