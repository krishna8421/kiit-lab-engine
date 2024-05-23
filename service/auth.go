package service

import (
	"kiit-lab-engine/core/db"
)

type AuthService interface {
	Register() error
	Login() error
}

type authService struct {
	db *db.DBClient
}

func NewAuthService(db *db.DBClient) AuthService {
	return &authService{
		db: db,
	}
}

func (a *authService) Register() error {
	// Use a.db to interact with the database...
	// Implement the register logic
	return nil
}

func (a *authService) Login() error {
	// Use a.db to interact with the database...
	// Implement the login logic
	return nil
}
