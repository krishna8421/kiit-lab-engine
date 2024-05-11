package service

type AuthService interface {
	Register()
	Login()
}

// implement the AuthService interface
type authService struct {
}

func NewAuthService() AuthService {
	return &authService{}
}

func (a *authService) Register() {
	// implement the register logic
}

func (a *authService) Login() {
	// implement the login logic
}
