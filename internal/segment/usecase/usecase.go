package usecase

import (
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/segment/delivery"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/segment/usecase/dto"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/segment/usecase/errors"
	"log/slog"
)

type segmentUseCase struct {
	repo   Repository
	logger *slog.Logger
}

func NewUseCase(repo Repository, logger *slog.Logger) delivery.UseCase {
	return &segmentUseCase{
		repo:   repo,
		logger: logger.WithGroup("segmentUseCase"),
	}
}

func (u *segmentUseCase) AddSegment(segmentData dto.SegmentDTO) error {
	err := u.repo.AddSegment(
		segmentData.Name,
		dto.ParseUserPercentage(segmentData.UserPercentage),
		dto.ParseTTL(segmentData.TTL),
	)
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
