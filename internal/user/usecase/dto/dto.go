package dto

import "time"

type UserDTO struct {
	UserID int `json:"user_id"`
}

type UserInputDTO UserDTO

type UserSegmentsInputDTO struct {
	UserID          int            `json:"user_id"`
	OldSegmentNames []string       `json:"old_segment_names,omitempty"`
	NewSegmentNames []string       `json:"new_segment_names,omitempty"`
	TTL             *time.Duration `json:"ttl,omitempty"`
}

type UserSegmentsOutputDTO struct {
	SegmentNames []string `json:"segment_names,omitempty"`
}

const YearMonthLayout = "2006-01"

type UserHistoryInputDTO struct {
	YearMonth string `json:"year_month"`
}
