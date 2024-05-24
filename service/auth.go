package service

import (
	"context"
	"fmt"

	"kiit-lab-engine/db"
	"kiit-lab-engine/lib/jwt"
	"kiit-lab-engine/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(ctx context.Context, input RegisterInput) (*db.UserModel, error)
	Login(ctx context.Context, input LoginInput) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetRefreshTokenDuration() int64
	GetAccessTokenDuration() int64
}

type authService struct {
	userRepo   repository.UserRepository
	jwtManager *jwt.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *jwt.JWTManager) AuthService {
	return &authService{
		userRepo:   userRepo,
		jwtManager: jwtManager,
	}
}

type RegisterInput struct {
	Name     string
	Email    string
	Password string
}

func (a *authService) Register(ctx context.Context, input RegisterInput) (*db.UserModel, error) {
	user, err := a.userRepo.CreateNewUser(ctx, repository.NewUserInput{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return nil, err
	}
	return user, nil
}

type LoginInput struct {
	Email    string
	Password string
}

func (a *authService) Login(ctx context.Context, input LoginInput) (string, string, error) {
	user, err := a.userRepo.GetUserFromEmail(ctx, input.Email)
	if err != nil {
		return "", "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		return "", "", fmt.Errorf("invalid credentials")
	}

	accessToken, err := a.jwtManager.Generate(map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
	}, false)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := a.jwtManager.Generate(map[string]interface{}{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role,
	}, true)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (a *authService) RefreshToken(ctx context.Context, refreshToken string) (string, error) {
	newAccessToken, err := a.jwtManager.Refresh(refreshToken)
	if err != nil {
		return "", err
	}
	return newAccessToken, nil
}

func (a *authService) GetRefreshTokenDuration() int64 {
	return a.jwtManager.RefreshTokenDuration
}

func (a *authService) GetAccessTokenDuration() int64 {
	return a.jwtManager.AccessTokenDuration
}
