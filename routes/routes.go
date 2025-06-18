package routes

import (
	"jaga/config"
	"jaga/consts"
	"jaga/controllers"
	"jaga/middleware"
	"jaga/repositories"
	"jaga/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes() *gin.Engine {
	r := gin.Default()

	userRepositories := repositories.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepositories)
	userController := controllers.NewUserController(userService)
	authController := controllers.NewAuthController(userService)

	r.POST("/login", authController.Login)

	r.POST("/users", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), userController.CreateUser)

	return r
}
