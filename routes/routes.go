package routes

import (
	"kiit-lab-engine/controllers"
	"kiit-lab-engine/db"
	"kiit-lab-engine/lib/jwt"
	"kiit-lab-engine/repository"
	"kiit-lab-engine/service"

	"github.com/gin-gonic/gin"
)

func InitRoutes(r *gin.Engine, dbClient *db.DBClient, jwtManager *jwt.JWTManager) {
	userRepo := repository.NewUserRepository(dbClient)

	userService := service.NewUserService(userRepo)
	authService := service.NewAuthService(userRepo, jwtManager)

	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)

	auth := r.Group("/auth")
	auth.POST("/register", authController.Register)
	auth.POST("/login", authController.Login)

	user := r.Group("/user")
	user.GET("/:id", userController.GetUser)
}
