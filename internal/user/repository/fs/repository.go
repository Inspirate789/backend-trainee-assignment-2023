package fs

import (
	"encoding/csv"
	"fmt"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/models"
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/user/usecase"
	"go.uber.org/multierr"
	"log/slog"
	"os"
)

type fsUserRepository struct {
	volumePath string
	logger     *slog.Logger
}

func NewFsRepository(volumePath string, logger *slog.Logger) usecase.FsRepository {
	return &fsUserRepository{
		volumePath: volumePath,
		logger:     logger.WithGroup("fsUserRepository"),
	}
}

func (r *fsUserRepository) SaveUserHistory(events []models.SegmentEvent, reportID string) (filename string, err error) {
	filename = fmt.Sprintf("%s/%s.csv", r.volumePath, reportID)
	file, err := os.Create(filename)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err = multierr.Combine(err, file.Close())
	}(file)

	w := csv.NewWriter(file)
	defer w.Flush()

	for _, event := range events {
		err = w.Write(event.ToCsvStrings())
		if err != nil {
			return "", err
		}
	}

	return filename, nil
}
