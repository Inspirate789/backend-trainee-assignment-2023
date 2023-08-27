package delivery

import "github.com/Inspirate789/backend-trainee-assignment-2023/internal/user/usecase/dto"

type UseCase interface {
	AddUser(userData dto.UserDTO) error
	RemoveUser(userData dto.UserInputDTO) error
	ChangeUserSegments(userData dto.UserSegmentsInputDTO) error
	GetUserSegments(userData dto.UserInputDTO) (dto.UserSegmentsOutputDTO, error)
	SaveUserHistory(userData dto.UserHistoryInputDTO, reportID string) (string, error)
}
