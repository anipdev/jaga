package controllers

import (
	"math"
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

func (ctrl *UserController) GetUsers(c *gin.Context) {
	var req dto.GetUsersRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid query parameters: " + err.Error()})
		return
	}

	usersModel, totalItems, err := ctrl.UserService.GetUsers(req.Page, req.ItemsPerPage, req.SortBy, req.SortDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve users: " + err.Error()})
		return
	}

	userDTOs := make([]dto.UserDTO, len(usersModel))
	for i, user := range usersModel {
		userDTOs[i] = dto.UserDTO{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		}
	}

	totalPages := 0
	if req.ItemsPerPage > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(req.ItemsPerPage)))
	}

	c.JSON(http.StatusOK, dto.GetUsersResponse{
		Message:      "Users retrieved successfully",
		Users:        userDTOs,
		TotalItems:   int(totalItems),
		Page:         req.Page,
		ItemsPerPage: req.ItemsPerPage,
		TotalPages:   totalPages,
	})
}

func (ctrl *UserController) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	userModel, err := ctrl.UserService.GetUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, dto.GetUserByIDResponse{
				Message: "User not found",
				User:    dto.UserDTO{},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.GetUserByIDResponse{
			Message: "Failed to retrieve user: " + err.Error(),
			User:    dto.UserDTO{},
		})
		return
	}

	userDTO := dto.UserDTO{
		ID:    userModel.ID,
		Name:  userModel.Name,
		Email: userModel.Email,
		Role:  userModel.Role,
	}

	c.JSON(http.StatusOK, dto.GetUserByIDResponse{
		Message: "User retrieved successfully",
		User:    userDTO,
	})
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

func (ctrl *UserController) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var req dto.UpdateUserRequest

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

	existingUser, err := ctrl.UserService.GetUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get existing user: " + err.Error()})
		return
	}

	if creatorRoleStr == "admin" && existingUser.Role == "admin" && userID != c.GetString("user_id") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin cannot update another admin"})
		return
	}

	updatedUser := &models.User{
		ID:           userID,
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: req.Password,
		Role:         req.Role,
	}

	err = ctrl.UserService.UpdateUser(updatedUser, creatorRoleStr)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err.Error() == "admin cannot update a user to 'admin' role" {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.UpdateUserResponse{
		Message: "User updated successfully!",
	})
}

func (ctrl *UserController) DeleteUser(c *gin.Context) {
	userID := c.Param("id")

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

	existingUser, err := ctrl.UserService.GetUserByID(userID)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get existing user: " + err.Error()})
		return
	}

	if creatorRoleStr == "admin" && existingUser.Role == "admin" && userID != c.GetString("user_id") {
		c.JSON(http.StatusForbidden, gin.H{"error": "Admin cannot delete another admin"})
		return
	}

	err = ctrl.UserService.DeleteUser(userID, creatorRoleStr)
	if err != nil {
		if err.Error() == "user not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.DeleteUserResponse{
		Message: "User deleted successfully!",
	})

}
