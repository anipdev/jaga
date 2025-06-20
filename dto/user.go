package dto

type UserDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
type GetUserByIDRequest struct {
	ID string `json:"id" binding:"required"`
}

type GetUserByIDResponse struct {
	Message string  `json:"message"`
	User    UserDTO `json:"user"`
}

type GetUsersRequest struct {
	Page         int    `json:"page" binding:"omitempty,min=1"`
	ItemsPerPage int    `json:"items_per_page" binding:"omitempty,min=1"`
	SortBy       string `json:"sort_by" binding:"omitempty,oneof=id name email role created_at updated_at"`
	SortDir      string `json:"sort_dir" binding:"omitempty,oneof=asc desc"`
}

type GetUsersResponse struct {
	Message      string    `json:"message"`
	Users        []UserDTO `json:"users"`
	TotalItems   int       `json:"total_items"`
	Page         int       `json:"page"`
	ItemsPerPage int       `json:"items_per_page"`
	TotalPages   int       `json:"total_pages"`
}

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Role     string `json:"role" binding:"required,oneof=super_user admin technician manager"`
}

type CreateUserResponse struct {
	Message string `json:"message"`
}

type UpdateUserRequest struct {
	ID       string `json:"id" binding:"required"`
	Name     string `json:"name,omitempty"`
	Email    string `json:"email,omitempty" binding:"omitempty,email"`
	Password string `json:"password,omitempty" binding:"omitempty,min=8"`
	Role     string `json:"role,omitempty" binding:"omitempty,oneof=super_user admin technician manager"`
}

type UpdateUserResponse struct {
	Message string  `json:"message"`
	User    UserDTO `json:"user"`
}

type DeleteUserRequest struct {
	ID string `json:"id" binding:"required"`
}

type DeleteUserResponse struct {
	Message string `json:"message"`
	UserID  string `json:"userId"`
}
