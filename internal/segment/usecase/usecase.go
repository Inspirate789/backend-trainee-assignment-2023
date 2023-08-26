package usecase

import (
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/segment/delivery"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/segment/usecase/dto"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/segment/usecase/errors"
	"log/slog"
	"time"
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

func (u *segmentUseCase) parseUserPercentage(segmentData dto.SegmentDTO) float64 {
	if segmentData.UserPercentage == nil {
		return EmptySegment
	}
	return *segmentData.UserPercentage
}

func (u *segmentUseCase) parseTTL(segmentData dto.SegmentDTO) time.Duration {
	if segmentData.TTL == nil {
		return NoTTL
	}
	return *segmentData.TTL
}

func (u *segmentUseCase) AddSegment(segmentData dto.SegmentDTO) error {
	err := u.repo.AddSegment(segmentData.Name, u.parseUserPercentage(segmentData), u.parseTTL(segmentData))

	if err != nil {
		u.logger.Error(err.Error())
		return errors.AddSegmentErr
	}
	u.logger.Info("new segment added", slog.String("name", segmentData.Name))

	return nil
}

func (u *segmentUseCase) RemoveSegment(segmentData dto.SegmentDTO) error {
	err := u.repo.RemoveSegment(segmentData.Name)

	if err != nil {
		u.logger.Error(err.Error())
		return errors.RemoveSegmentErr
	}
	u.logger.Info("segment removed", slog.String("name", segmentData.Name))

	return nil
}
