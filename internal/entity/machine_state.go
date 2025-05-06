package entity

type StateMachine struct {
	transitions map[Status][]Status
}

func NewStateMachine() *StateMachine {
	transitions := map[Status][]Status{
		StatusRegistered: {
			StatusOnLap,
			StatusDisqualified,
		},
		StatusOnLap: {
			StatusOnFiringRange,
			StatusNotFinished,
			StatusFinished,
			StatusOnLap,
			StatusOnPenaltyLap, // for case when penalty laps > 1
		},
		StatusOnFiringRange: {
			StatusOnPenaltyLap,
			StatusOnLap,
		},
		StatusOnPenaltyLap: {
			StatusOnLap,
			StatusOnPenaltyLap,
		},
	}

	return &StateMachine{
		transitions: transitions,
	}
}

func (sm *StateMachine) MoveTo(current, next Status) bool {
	allowedStates, exists := sm.transitions[current]
	if !exists {
		return false
	}

	for _, state := range allowedStates {
		if state == next {
			return true
		}
	}
	return false
}

func (sm *StateMachine) IsTerminal(status Status) bool {
	_, exists := sm.transitions[status]
	if !exists {
		return true
	}

	return false
}
