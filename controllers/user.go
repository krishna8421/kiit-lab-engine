package controllers

import (
	"github.com/gin-gonic/gin"
	"kiit-lab-engine/service"
	"net/http"
)

type UserController struct {
	userService service.UserService
}

func NewUserController(userService service.UserService) *UserController {
	return &UserController{userService: userService}
}

func (uc *UserController) GetUser(c *gin.Context) {
	// Implement the get user logic using the service
	id := c.Param("id")
	user, err := uc.userService.GetUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}
