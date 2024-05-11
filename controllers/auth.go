package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "login",
	})
}

func Register(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "register",
	})
}
