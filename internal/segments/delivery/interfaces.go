package delivery

type SegmentUseCase interface {
	AddSegment(name string) error
	RemoveSegment(name string) error
}
