package delivery

import "github.com/Inspirate789/backend-trainee-assignment-2023/internal/segment/usecase/dto"

type SegmentUseCase interface {
	AddSegment(segmentData dto.SegmentDTO) error
	RemoveSegment(segmentData dto.SegmentDTO) error
}