package delivery

import "github.com/Inspirate789/backend-trainee-assignment-2023/internal/user/usecase/dto"

type UseCase interface {
	AddUser(userData dto.UserDTO) error
	RemoveUser(userID int) error
	ChangeUserSegments(userData dto.UserSegmentsInputDTO) error
	GetUserSegments(userID int) (dto.UserSegmentsOutputDTO, error)
	SaveUserHistory(yearMonth, reportID string) (string, error)
}
