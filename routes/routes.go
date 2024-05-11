package routes

import (
	"kiit-lab-engine/controllers"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine) {
	auth := r.Group("/auth")
	auth.POST("/register", controllers.Register)
	auth.POST("/login", controllers.Login)

	user := r.Group("/user")
	user.GET("/:id", controllers.GetUser)

}
