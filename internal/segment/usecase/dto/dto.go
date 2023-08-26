package dto

import "time"

type SegmentDTO struct {
	Name           string         `json:"name"`
	UserPercentage *float64       `json:"user_percentage,omitempty"`
	TTL            *time.Duration `json:"ttl,omitempty"`
}
