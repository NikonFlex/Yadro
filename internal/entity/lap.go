package entity

import (
	"time"
)

type Lap struct {
	StartTime time.Time
	EndTime   time.Time
}
