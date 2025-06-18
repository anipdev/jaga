package controllers

import (
	"jaga/dto"
	"jaga/services"
	"jaga/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuthController struct {
	UserService services.UserService
}

func NewAuthController(userService services.UserService) *AuthController {
	return &AuthController{UserService: userService}
}

func (ctrl *AuthController) Login(c *gin.Context) {
	var req dto.LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Bad Request",
		})
		return
	}

	user, err := ctrl.UserService.GetUserByEmail(req.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid credentials"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to retrieve user",
		})
		return
	}

	if err := utils.ComparePassword(user.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token",
		})
		return
	}

	c.JSON(http.StatusOK, dto.LoginResponse{
		Token: token,
	})
}
