package usecase

import (
	"github.com/Inspirate789/backend-trainee-assignment-2023/internal/models"
	"time"
)

type SqlRepository interface {
	AddUser(userID int) error
	RemoveUser(userID int) error
	ChangeUserSegments(userID int, oldSegmentNames, newSegmentNames []string, ttl time.Duration) error
	GetUserSegments(userID int) ([]string, error)
	GetUserHistory(year, month int) ([]models.SegmentEvent, error)
}

type FsRepository interface {
	SaveUserHistory(events []models.SegmentEvent, reportID string) (string, error)
}
