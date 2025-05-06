package entity

import "errors"

type Status int

const (
	StatusUnknown Status = iota
	StatusRegistered
	StatusOnLap
	StatusOnPenaltyLap
	StatusOnFiringRange
	StatusNotFinished
	StatusDisqualified
	StatusFinished
)

var ErrStateIsNotAccepted = errors.New("state is not accepted")
