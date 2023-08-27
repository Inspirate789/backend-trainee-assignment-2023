package errors

type UserUseCaseError string

func (e UserUseCaseError) Error() string {
	return string(e)
}

const (
	AddUserErr            UserUseCaseError = "Cannot add new user"
	RemoveUserErr         UserUseCaseError = "Cannot remove user"
	ChangeUserSegmentsErr UserUseCaseError = "Cannot change user segments"
	GetUserSegmentsErr    UserUseCaseError = "Cannot get user segments"
	ParseDateErr          UserUseCaseError = "Cannot parse date"
	GetUserHistoryErr     UserUseCaseError = "Cannot get users history"
	SaveUserHistoryErr    UserUseCaseError = "Cannot save users history to file"
)
