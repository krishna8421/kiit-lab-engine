package controllers

import (
	"kiit-lab-engine/service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

func (a *AuthController) Register(c *gin.Context) {
	var input service.RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := a.authService.Register(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (a *AuthController) Login(c *gin.Context) {
	var input service.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := a.authService.Login(c.Request.Context(), input)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Get the token durations from the AuthService
	accessTokenDuration := a.authService.GetAccessTokenDuration()
	refreshTokenDuration := a.authService.GetRefreshTokenDuration()

	// Convert the token durations from seconds to durations
	accessTokenDurationInTime := time.Duration(accessTokenDuration) * time.Second
	refreshTokenDurationInTime := time.Duration(refreshTokenDuration) * time.Second

	// Set a secure HTTP-only cookie with the access token
	c.SetCookie("access_token", accessToken, int(accessTokenDurationInTime.Seconds()), "/", "", true, true)
	c.Writer.Header().Add("Set-Cookie", "access_token="+accessToken+"; Path=/; HttpOnly; Secure; SameSite=Strict")

	// Set a secure HTTP-only cookie with the refresh token
	c.SetCookie("refresh_token", refreshToken, int(refreshTokenDurationInTime.Seconds()), "/", "", true, true)
	c.Writer.Header().Add("Set-Cookie", "refresh_token="+refreshToken+"; Path=/; HttpOnly; Secure; SameSite=Strict")

	c.JSON(http.StatusOK, gin.H{"message": "Logged in successfully"})
}
