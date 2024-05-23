package routes

import (
	"kiit-lab-engine/controllers"
	"kiit-lab-engine/core/db"
	"kiit-lab-engine/service"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, dbClient *db.DBClient) {
	authService := service.NewAuthService(dbClient)
	userService := service.NewUserService(dbClient)

	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)

	auth := r.Group("/auth")
	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)

	user := r.Group("/user")
	user.GET("/:id", userController.GetUser)
}
