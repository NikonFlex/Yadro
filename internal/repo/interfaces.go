package repo

import (
	"Solution/internal/entity"
	"go.uber.org/zap"
	"time"
)

type LapRepository interface {
	StartLap(competitorId int, startTime time.Time) error
	FinishLap(competitorId int, finishTime time.Time) error
	StartPenaltyLap(competitorId int, startTime time.Time) error
	FinishPenaltyLap(competitorId int, finishTime time.Time) error
}

type FiringRangeRepository interface {
	StartFiringRange(competitorId int) error
	ShotTaken(competitorId int, result bool) error
}

type CompetitorRepository interface {
	Get(id int) (entity.Competitor, error)
	GetAll() []entity.Competitor
	Register(id int) error
	SetStartTime(id int, startTime time.Time) error
	Start(id int, startTime time.Time) error
	Finish(id int, finishTime time.Time) error
	CantContinue(id int, comment string) error
	GetStatus(competitorId int) (entity.Status, error)
	SetStatus(competitorId int, status entity.Status) error
}

var _ CompetitorRepository = (*repository)(nil)
var _ LapRepository = (*repository)(nil)
var _ FiringRangeRepository = (*repository)(nil)

type repository struct {
	logger      *zap.Logger
	competitors map[int]*entity.Competitor
}

func New(logger *zap.Logger) *repository {
	return &repository{
		logger:      logger,
		competitors: map[int]*entity.Competitor{},
	}
}
