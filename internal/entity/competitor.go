package entity

import (
	"errors"
	"time"
)

type Competitor struct {
	ID             int
	ScheduledStart time.Time
	Start          time.Time
	Finish         time.Time
	Laps           []Lap
	PenaltyLaps    []Lap
	FiringResults  []FiringResult
	Status         Status
	Comment        string
}

var ErrCompetitorNotFound = errors.New("competitor not found")
var ErrCompetitorAlreadyExists = errors.New("competitor already exists")
