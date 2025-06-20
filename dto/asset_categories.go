package dto

type AssetCategoryDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type CreateAssetCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

type CreateAssetCategoryResponse struct {
	Message string `json:"message"`
}

type GetAssetCategoryByIDResponse struct {
	Message       string           `json:"message"`
	AssetCategory AssetCategoryDTO `json:"asset_category"`
}

type GetAssetCategoriesRequest struct {
	Page         int    `json:"page" binding:"omitempty,min=1"`
	ItemsPerPage int    `json:"items_per_page" binding:"omitempty,min=1"`
	SortBy       string `json:"sort_by" binding:"omitempty,oneof=id name created_at updated_at"`
	SortDir      string `json:"sort_dir" binding:"omitempty,oneof=asc desc"`
}

type GetAssetCategoriesResponse struct {
	Message         string             `json:"message"`
	AssetCategories []AssetCategoryDTO `json:"asset_categories"`
}

type UpdateAssetCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2,max=100"`
}

type UpdateAssetCategoryResponse struct {
	Message string `json:"message"`
}

type DeleteAssetCategoryResponse struct {
	Message string `json:"message"`
}
