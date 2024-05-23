package controllers

import (
	"github.com/gin-gonic/gin"
	"kiit-lab-engine/service"
	"net/http"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (ac *AuthController) Login(c *gin.Context) {
	// Implement the login logic using the service
	ac.authService.Login()
	c.JSON(http.StatusOK, gin.H{"message": "login"})
}

func (ac *AuthController) Register(c *gin.Context) {
	// Implement the register logic using the service
	ac.authService.Register()
	c.JSON(http.StatusOK, gin.H{"message": "register"})
}
