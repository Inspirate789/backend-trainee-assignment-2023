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

func ParseTTL(ptr *time.Duration) time.Duration {
	if ptr == nil {
		return NoTTL
	}
	return *ptr
}

type SegmentDTO struct {
	Name           string         `json:"name"`
	UserPercentage *float64       `json:"user_percentage,omitempty"`
	TTL            *time.Duration `json:"ttl,omitempty"`
}
