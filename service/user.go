package service

type UserService interface {
	GetUser()
}

// implement the UserService interface
type userService struct {
}

func NewUserService() UserService {
	return &userService{}
}

func (u *userService) GetUser() {
	// implement the get user logic
}
