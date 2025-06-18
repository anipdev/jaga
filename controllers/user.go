package controllers

import (
	"net/http"

	"jaga/dto"
	"jaga/models"
	"jaga/services"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserService services.UserService
}

func NewUserController(userService services.UserService) *UserController {
	return &UserController{UserService: userService}
}

func (ctrl *UserController) CreateUser(c *gin.Context) {
	var req dto.CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	creatorRole, exists := c.Get("role")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Creator role not found in context"})
		return
	}
	creatorRoleStr, ok := creatorRole.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid creator role type"})
		return
	}

	newUser := &models.User{
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
		Role:         req.Role,
	}

	_, err := ctrl.UserService.CreateUser(newUser, creatorRoleStr)
	if err != nil {
		if err.Error() == "admin cannot create a user with 'admin' role" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "email already registered" {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.CreateUserResponse{
		Message: "User created successfully!",
	})
}
