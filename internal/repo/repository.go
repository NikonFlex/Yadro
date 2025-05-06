package repo

import (
	"Solution/internal/entity"
	"go.uber.org/zap"
	"time"
)

func (r repository) Get(id int) (entity.Competitor, error) {
	r.logger.Debug("Getting competitor", zap.Int("id", id))
	competitor, ok := r.competitors[id]
	if !ok {
		r.logger.Error("Competitor not found", zap.Int("id", id))
		return entity.Competitor{}, entity.ErrCompetitorNotFound
	}

	return *competitor, nil
}

func (r repository) GetAll() []entity.Competitor {
	result := make([]entity.Competitor, 0, len(r.competitors))
	for _, competitor := range r.competitors {
		result = append(result, *competitor)
	}
	return result
}

func (r repository) Register(id int) error {
	r.logger.Debug("Registering competitor", zap.Int("id", id))
	if _, ok := r.competitors[id]; ok {
		r.logger.Error("Competitor already exists", zap.Int("id", id))
		return entity.ErrCompetitorAlreadyExists
	}

	competitor := &entity.Competitor{
		ID: id,
	}

	r.competitors[id] = competitor
	return nil
}

func (r repository) SetStartTime(id int, startTime time.Time) error {
	r.logger.Debug("Setting start time", zap.Int("id", id))
	competitor, ok := r.competitors[id]
	if !ok {
		r.logger.Error("Competitor not found", zap.Int("id", id))
		return entity.ErrCompetitorNotFound
	}

	competitor.ScheduledStart = startTime
	return nil
}

func (r repository) Start(id int, startTime time.Time) error {
	r.logger.Debug("Starting competitor", zap.Int("id", id))
	competitor, ok := r.competitors[id]
	if !ok {
		r.logger.Error("Competitor not found", zap.Int("id", id))
		return entity.ErrCompetitorNotFound
	}

	competitor.Start = startTime
	return nil
}

func (r repository) Finish(id int, finishTime time.Time) error {
	r.logger.Debug("Finishing competitor", zap.Int("id", id))
	competitor, ok := r.competitors[id]
	if !ok {
		r.logger.Error("Competitor not found", zap.Int("id", id))
		return entity.ErrCompetitorNotFound
	}

	competitor.Finish = finishTime
	return nil
}

func (r repository) CantContinue(id int, comment string) error {
	r.logger.Debug("Set cant continue to competitor", zap.Int("id", id))
	_, ok := r.competitors[id]
	if !ok {
		r.logger.Error("Competitor not found", zap.Int("id", id))
		return entity.ErrCompetitorNotFound
	}

	r.competitors[id].Comment = comment
	return nil
}

func (r repository) GetStatus(competitorId int) (entity.Status, error) {
	r.logger.Debug("Getting competitor status", zap.Int("id", competitorId))
	competitor, ok := r.competitors[competitorId]
	if !ok {
		r.logger.Error("Competitor not found", zap.Int("id", competitorId))
		return entity.StatusUnknown, entity.ErrCompetitorNotFound
	}

	return competitor.Status, nil
}

func (r repository) SetStatus(competitorId int, status entity.Status) error {
	r.logger.Debug("Setting competitor status", zap.Int("id", competitorId))
	competitor, ok := r.competitors[competitorId]
	if !ok {
		r.logger.Error("Competitor not found", zap.Int("id", competitorId))
		return entity.ErrCompetitorNotFound
	}

	competitor.Status = status
	return nil
}

func (r repository) StartLap(competitorId int, startTime time.Time) error {
	r.logger.Debug("Starting competitor lap", zap.Int("id", competitorId))
	competitor, ok := r.competitors[competitorId]
	if !ok {
		r.logger.Error("Competitor not found", zap.Int("id", competitorId))
		return entity.ErrCompetitorNotFound
	}

	lap := entity.Lap{
		StartTime: startTime,
	}
	competitor.Laps = append(competitor.Laps, lap)
	return nil
}

func (r repository) FinishLap(competitorId int, finishTime time.Time) error {
	r.logger.Debug("Finishing competitor lap", zap.Int("id", competitorId))
	competitor, ok := r.competitors[competitorId]
	if !ok {
		r.logger.Error("Competitor not found", zap.Int("id", competitorId))
		return entity.ErrCompetitorNotFound
	}

	competitor.Laps[len(competitor.Laps)-1].EndTime = finishTime
	return nil
}

func (r repository) StartPenaltyLap(competitorId int, startTime time.Time) error {
	r.logger.Debug("Starting competitor penalty lap", zap.Int("id", competitorId))
	competitor, ok := r.competitors[competitorId]
	if !ok {
		r.logger.Error("Competitor not found", zap.Int("id", competitorId))
		return entity.ErrCompetitorNotFound
	}

	lap := entity.Lap{
		StartTime: startTime,
	}
	competitor.PenaltyLaps = append(competitor.PenaltyLaps, lap)
	return nil
}

func (r repository) FinishPenaltyLap(competitorId int, finishTime time.Time) error {
	r.logger.Debug("Finishing competitor penalty lap", zap.Int("id", competitorId))
	competitor, ok := r.competitors[competitorId]
	if !ok {
		r.logger.Error("Competitor not found", zap.Int("id", competitorId))
		return entity.ErrCompetitorNotFound
	}

	competitor.PenaltyLaps[len(competitor.PenaltyLaps)-1].EndTime = finishTime
	return nil
}

func (r repository) StartFiringRange(competitorId int) error {
	r.logger.Debug("Starting competitor firing range", zap.Int("id", competitorId))
	competitor, ok := r.competitors[competitorId]
	if !ok {
		r.logger.Error("Competitor not found", zap.Int("id", competitorId))
		return entity.ErrCompetitorNotFound
	}

	competitor.FiringResults = append(competitor.FiringResults, entity.FiringResult{})
	return nil
}

func (r repository) ShotTaken(competitorId int, result bool) error {
	r.logger.Debug("Competitor took a shot", zap.Int("id", competitorId))
	competitor, ok := r.competitors[competitorId]
	if !ok {
		r.logger.Error("Competitor not found", zap.Int("id", competitorId))
		return entity.ErrCompetitorNotFound
	}

	if result {
		competitor.FiringResults[len(competitor.FiringResults)-1].Hits++
	}

	return nil
}
