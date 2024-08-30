package models

import "time"

type Reservation struct {
	ID        int64
	RoomID    string
	StartTime time.Time
	EndTime   time.Time
}
