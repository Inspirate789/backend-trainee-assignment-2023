package usecase

import "time"

type Repository interface {
	AddSegment(name string, userPercentage float64, ttl time.Duration) error
	RemoveSegment(name string) error
}
