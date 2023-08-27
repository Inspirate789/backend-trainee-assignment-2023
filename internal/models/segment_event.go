package models

import (
	"strconv"
	"time"
)

type SegmentEvent struct {
	UserID           int       `db:"user_id"`
	SegmentName      string    `db:"segment_name"`
	Operation        string    `db:"operation"`
	RegistrationDate time.Time `db:"registration_date"`
}

func (se *SegmentEvent) ToCsvStrings() []string {
	return []string{
		strconv.FormatInt(int64(se.UserID), 10),
		se.SegmentName,
		se.Operation,
		se.RegistrationDate.Format(time.DateTime),
	}
}
