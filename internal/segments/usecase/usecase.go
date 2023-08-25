package usecase

import (
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/segments/delivery"
	"log/slog"
)

type segmentUseCase struct {
	repo   SegmentRepository
	logger *slog.Logger
}

func NewSegmentUseCase(repo SegmentRepository, logger *slog.Logger) delivery.SegmentUseCase {
	return &segmentUseCase{
		repo:   repo,
		logger: logger.WithGroup("segmentUseCase"),
	}
}

func (s *segmentUseCase) AddSegment(name string) error {
	//TODO implement me
	panic("implement me")
}

func (s *segmentUseCase) RemoveSegment(name string) error {
	//TODO implement me
	panic("implement me")
}
