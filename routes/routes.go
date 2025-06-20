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

	userRoutes := r.Group("/users")
	{
		userRoutes.POST("", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), userController.CreateUser)
		userRoutes.GET("", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), userController.GetUsers)
		userRoutes.GET("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), userController.GetUserByID)
		userRoutes.PUT("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), userController.UpdateUser)
		userRoutes.DELETE("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), userController.DeleteUser)
	}

	assetCategoryRepository := repositories.NewAssetCategoryRepository(config.DB)
	assetCategoryService := services.NewAssetCategoryService(assetCategoryRepository)
	assetCategoryController := controllers.NewAssetCategoryController(assetCategoryService)

	assetCategoryRoutes := r.Group("/asset-categories")
	{
		assetCategoryRoutes.POST("", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), assetCategoryController.CreateAssetCategory)
		assetCategoryRoutes.GET("", middleware.RequireRole(consts.AllRoles...), assetCategoryController.GetAssetCategories)
		assetCategoryRoutes.GET("/:id", middleware.RequireRole(consts.AllRoles...), assetCategoryController.GetAssetCategoryByID)
		assetCategoryRoutes.PUT("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), assetCategoryController.UpdateAssetCategory)
		assetCategoryRoutes.DELETE("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), assetCategoryController.DeleteAssetCategory)
	}

	return r
}
