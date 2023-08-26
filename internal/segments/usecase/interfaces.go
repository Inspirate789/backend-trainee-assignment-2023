package usecase

import "time"

const (
	EmptySegment float64       = 0
	NoTTL        time.Duration = 0
)

type SegmentRepository interface {
	AddSegment(name string, userPercentage float64, ttl time.Duration) error
	RemoveSegment(name string) error
}
