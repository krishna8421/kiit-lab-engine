package service

import (
	"context"
	"kiit-lab-engine/db"
	"kiit-lab-engine/repository"
)

type UserService interface {
	GetUserByID(ctx context.Context, id string) (*db.UserModel, error)
	GetUserByEmail(ctx context.Context, email string) (*db.UserModel, error)
	CreateNewUser(ctx context.Context, user repository.NewUserInput) (*db.UserModel, error)
	UpdateUser(ctx context.Context, id string, updateUser repository.UpdateUserInput) (*db.UserModel, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (u *userService) GetUserByID(ctx context.Context, id string) (*db.UserModel, error) {
	return u.repo.GetUserFromId(ctx, id)
}

func (u *userService) GetUserByEmail(ctx context.Context, email string) (*db.UserModel, error) {
	return u.repo.GetUserFromEmail(ctx, email)
}

func (u *userService) CreateNewUser(ctx context.Context, user repository.NewUserInput) (*db.UserModel, error) {
	return u.repo.CreateNewUser(ctx, user)
}

func (u *userService) UpdateUser(ctx context.Context, id string, updateUser repository.UpdateUserInput) (*db.UserModel, error) {
	return u.repo.UpdateUser(ctx, id, updateUser)
}
