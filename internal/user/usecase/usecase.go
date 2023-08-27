package usecase

import (
	segmentDTO "github.com/Inspirate789/backend-trainee-assignment-2023/internal/segment/usecase/dto"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/user/delivery"
	userDTO "github.com/Inspirate789/backend-trainee-assignment-2023/internal/user/usecase/dto"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/user/usecase/errors"
	"log/slog"
	"time"
)

type userUseCase struct {
	sqlRepo SqlRepository
	fsRepo  FsRepository
	logger  *slog.Logger
}

func NewUseCase(sqlRepo SqlRepository, fsRepo FsRepository, logger *slog.Logger) delivery.UseCase {
	return &userUseCase{
		sqlRepo: sqlRepo,
		fsRepo:  fsRepo,
		logger:  logger.WithGroup("userUseCase"),
	}
}

func (u *userUseCase) AddUser(userData userDTO.UserDTO) error {
	err := u.sqlRepo.AddUser(userData.UserID)
	if err != nil {
		u.logger.Error(err.Error())
		return errors.AddUserErr
	}
	u.logger.Info("new user added", slog.Int("user_id", userData.UserID))

	return nil
}

func (u *userUseCase) RemoveUser(userData userDTO.UserInputDTO) error {
	err := u.sqlRepo.AddUser(userData.UserID)
	if err != nil {
		u.logger.Error(err.Error())
		return errors.RemoveUserErr
	}
	u.logger.Info("user removed", slog.Int("user_id", userData.UserID))

	return nil
}

func (u *userUseCase) ChangeUserSegments(userData userDTO.UserSegmentsInputDTO) error {
	err := u.sqlRepo.ChangeUserSegments(
		userData.UserID,
		userData.OldSegmentNames,
		userData.NewSegmentNames,
		segmentDTO.ParseTTL(userData.TTL),
	)
	if err != nil {
		u.logger.Error(err.Error())
		return errors.ChangeUserSegmentsErr
	}
	u.logger.Info(
		"user segments changed",
		slog.Int("user_id", userData.UserID),
		"old_segment_names", userData.OldSegmentNames,
		"new_segment_names", userData.NewSegmentNames,
	)

	return nil
}

func (u *userUseCase) GetUserSegments(userData userDTO.UserInputDTO) (userDTO.UserSegmentsOutputDTO, error) {
	segments, err := u.sqlRepo.GetUserSegments(userData.UserID)
	if err != nil {
		u.logger.Error(err.Error())
		return userDTO.UserSegmentsOutputDTO{}, errors.GetUserSegmentsErr
	}
	u.logger.Info("user segments received", slog.Int("user_id", userData.UserID))

	return userDTO.UserSegmentsOutputDTO{SegmentNames: segments}, nil
}

func (u *userUseCase) SaveUserHistory(userData userDTO.UserHistoryInputDTO, reportID string) (string, error) {
	date, err := time.Parse(userDTO.YearMonthLayout, userData.YearMonth)
	if err != nil {
		u.logger.Error(err.Error())
		return "", errors.ParseDateErr
	}
	year, month := date.Year(), int(date.Month())
	u.logger.Debug(
		"date extracted",
		slog.Int("year", year),
		slog.Int("month", month),
	)

	events, err := u.sqlRepo.GetUserHistory(year, month)
	if err != nil {
		u.logger.Error(err.Error())
		return "", errors.GetUserHistoryErr
	}

	filename, err := u.fsRepo.SaveUserHistory(events, reportID)
	if err != nil {
		u.logger.Error(err.Error())
		return "", errors.SaveUserHistoryErr
	}
	u.logger.Info(
		"report saved",
		slog.Int("year", year),
		slog.Int("month", month),
	)

	return filename, nil
}
