package routes

import (
	"jaga/config"
	"jaga/consts"
	"jaga/controllers"
	"jaga/middleware"
	"jaga/repositories"
	"jaga/services"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func RegisterRoutes() *gin.Engine {
	r := gin.Default()

	userRepositories := repositories.NewUserRepository(config.DB)
	userService := services.NewUserService(userRepositories)
	userController := controllers.NewUserController(userService)

	authController := controllers.NewAuthController(userService)

	assetCategoryRepository := repositories.NewAssetCategoryRepository(config.DB)
	assetCategoryService := services.NewAssetCategoryService(assetCategoryRepository)
	assetCategoryController := controllers.NewAssetCategoryController(assetCategoryService)

	assetRepository := repositories.NewAssetRepository(config.DB)
	assetService := services.NewAssetService(assetRepository, assetCategoryRepository)
	assetController := controllers.NewAssetController(assetService)

	maintenanceScheduleRepository := repositories.NewMaintenanceScheduleRepository(config.DB)
	maintenanceScheduleService := services.NewMaintenanceScheduleService(maintenanceScheduleRepository, assetRepository)
	maintenanceScheduleController := controllers.NewMaintenanceScheduleController(maintenanceScheduleService)

	maintenanceRecordRepository := repositories.NewMaintenanceRecordRepository(config.DB)
	maintenanceRecordService := services.NewMaintenanceRecordService(maintenanceRecordRepository, assetRepository, maintenanceScheduleRepository, userRepositories)
	maintenanceRecordController := controllers.NewMaintenanceRecordController(maintenanceRecordService)

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("/v1")
	{
		v1.POST("/login", authController.Login)

		userRoutes := v1.Group("/users")
		{
			userRoutes.POST("", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), userController.CreateUser)
			userRoutes.GET("", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), userController.GetUsers)
			userRoutes.GET("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), userController.GetUserByID)
			userRoutes.PUT("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), userController.UpdateUser)
			userRoutes.DELETE("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), userController.DeleteUser)
		}

		assetCategoryRoutes := v1.Group("/asset-categories")
		{
			assetCategoryRoutes.POST("", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), assetCategoryController.CreateAssetCategory)
			assetCategoryRoutes.GET("", middleware.RequireRole(consts.AllRoles...), assetCategoryController.GetAssetCategories)
			assetCategoryRoutes.GET("/:id", middleware.RequireRole(consts.AllRoles...), assetCategoryController.GetAssetCategoryByID)
			assetCategoryRoutes.PUT("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), assetCategoryController.UpdateAssetCategory)
			assetCategoryRoutes.DELETE("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), assetCategoryController.DeleteAssetCategory)
		}

		assetRoutes := v1.Group("/assets")
		{
			assetRoutes.POST("", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), assetController.CreateAsset)
			assetRoutes.GET("", middleware.RequireRole(consts.AllRoles...), assetController.GetAssets)
			assetRoutes.GET("/:id", middleware.RequireRole(consts.AllRoles...), assetController.GetAssetByID)
			assetRoutes.PUT("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), assetController.UpdateAsset)
			assetRoutes.DELETE("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), assetController.DeleteAsset)
		}

		maintenanceScheduleRoutes := v1.Group("/maintenance-schedules")
		{
			maintenanceScheduleRoutes.POST("", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), maintenanceScheduleController.CreateMaintenanceSchedule)
			maintenanceScheduleRoutes.GET("", middleware.RequireRole(consts.AllRoles...), maintenanceScheduleController.GetMaintenanceSchedules)
			maintenanceScheduleRoutes.GET("/:id", middleware.RequireRole(consts.AllRoles...), maintenanceScheduleController.GetMaintenanceScheduleByID)
			maintenanceScheduleRoutes.PUT("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), maintenanceScheduleController.UpdateMaintenanceSchedule)
			maintenanceScheduleRoutes.DELETE("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), maintenanceScheduleController.DeleteMaintenanceSchedule)
		}

		maintenanceRecordRoutes := v1.Group("/maintenance-records")
		{
			maintenanceRecordRoutes.POST("", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), maintenanceRecordController.CreateMaintenanceRecord)
			maintenanceRecordRoutes.GET("", middleware.RequireRole(consts.AllRoles...), maintenanceRecordController.GetMaintenanceRecords)
			maintenanceRecordRoutes.GET("/:id", middleware.RequireRole(consts.AllRoles...), maintenanceRecordController.GetMaintenanceRecordByID)
			maintenanceRecordRoutes.PUT("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), maintenanceRecordController.UpdateMaintenanceRecord)
			maintenanceRecordRoutes.PUT("/:id/status", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin, consts.RoleTechnician), maintenanceRecordController.UpdateMaintenanceRecordStatus)
			maintenanceRecordRoutes.DELETE("/:id", middleware.RequireRole(consts.RoleSuperUser, consts.RoleAdmin), maintenanceRecordController.DeleteMaintenanceRecord)
		}
	}

	return r
}
