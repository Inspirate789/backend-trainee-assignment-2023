package usecase

type SegmentRepository interface {
	AddSegment(name string) error
	RemoveSegment(name string) error
}
