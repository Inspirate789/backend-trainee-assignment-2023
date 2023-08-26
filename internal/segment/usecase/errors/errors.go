package errors

type SegmentUseCaseError string

func (e SegmentUseCaseError) Error() string {
	return string(e)
}

const (
	AddSegmentErr    SegmentUseCaseError = "Cannot add new segment"
	RemoveSegmentErr SegmentUseCaseError = "Cannot remove segment"
)
