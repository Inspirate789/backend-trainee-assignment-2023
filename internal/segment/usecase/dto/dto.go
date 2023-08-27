package dto

import "time"

const (
	EmptySegment float64       = 0
	NoTTL        time.Duration = 0
)

func ParseUserPercentage(ptr *float64) float64 {
	if ptr == nil {
		return EmptySegment
	}
	return *ptr
}

func ParseTTL(ptr *int) time.Duration {
	if ptr == nil {
		return NoTTL
	}
	return time.Duration(*ptr) * time.Hour
}

// SegmentDTO godoc
//
// swagger:model
type SegmentDTO struct {
	// Name
	// required: true
	// min length: 1
	// example: "AVITO_VOICE_MESSAGES"
	Name string `json:"name"`

	// UserPercentage - part of all users that segment contains (in %)
	// required: false
	// min: 0
	// max: 100
	// example: 50
	UserPercentage *float64 `json:"user_percentage,omitempty"`

	// TTL - segment existing time (in hours)
	// required: false
	// min: 1
	// example: 72
	TTL *int `json:"ttl,omitempty"`
}

// SegmentNameDTO godoc
//
// swagger:model
type SegmentNameDTO struct {
	// Name
	// required: true
	// min length: 1
	// example: "AVITO_VOICE_MESSAGES"
	Name string `json:"name"`
}
