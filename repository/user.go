package repository

import (
	"context"
	"fmt"
	"strings"

	"kiit-lab-engine/db"

	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	GetUserFromId(ctx context.Context, id string) (*db.UserModel, error)
	GetUserFromEmail(ctx context.Context, email string) (*db.UserModel, error)
	CreateNewUser(ctx context.Context, user NewUserInput) (*db.UserModel, error)
	UpdateUser(ctx context.Context, id string, updateUser UpdateUserInput) (*db.UserModel, error)
}

type userRepository struct {
	db *db.DBClient
}

func NewUserRepository(db *db.DBClient) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetUserFromId(ctx context.Context, id string) (*db.UserModel, error) {
	user, err := r.db.Prisma.User.FindUnique(
		db.User.ID.Equals(id),
	).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (r *userRepository) GetUserFromEmail(ctx context.Context, email string) (*db.UserModel, error) {
	user, err := r.db.Prisma.User.FindUnique(
		db.User.Email.Equals(email),
	).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

type NewUserInput struct {
	Name     string
	Email    string
	Password string
}

func (r *userRepository) CreateNewUser(ctx context.Context, user NewUserInput) (*db.UserModel, error) {
	if user.Name == "" || user.Email == "" || user.Password == "" {
		return nil, fmt.Errorf("invalid input: all fields are required")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	created, err := r.db.Prisma.User.CreateOne(
		db.User.Name.Set(user.Name),
		db.User.Email.Set(user.Email),
		db.User.Password.Set(string(hashedPassword)),
		db.User.Username.Set(strings.Split(user.Email, "@")[0]),
		db.User.Role.Set(db.RoleStudent),
	).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return created, nil
}

type UpdateUserInput struct {
	Name     *string
	Email    *string
	Password *string
}

func (r *userRepository) UpdateUser(ctx context.Context, id string, updateUser UpdateUserInput) (*db.UserModel, error) {
	updateFields := []db.UserSetParam{}

	if updateUser.Name != nil {
		updateFields = append(updateFields, db.User.Name.Set(*updateUser.Name))
	}
	if updateUser.Email != nil {
		updateFields = append(updateFields, db.User.Email.Set(*updateUser.Email))
	}
	if updateUser.Password != nil {
		updateFields = append(updateFields, db.User.Password.Set(*updateUser.Password))
	}

	updated, err := r.db.Prisma.User.FindUnique(
		db.User.ID.Equals(id),
	).Update(
		updateFields...,
	).Exec(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updated, nil
}
