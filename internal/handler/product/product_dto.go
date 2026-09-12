package product

type ProductCreateRequest struct {
	SKU      string `json:"sku"`
	Name     string `json:"name"`
	Unit     string `json:"unit"`
	IsActive int16  `json:"is_active"`
}

type ProductUpdateRequest struct {
	SKU      *string `json:"sku,omitempty"`
	Name     *string `json:"name,omitempty"`
	Unit     *string `json:"unit,omitempty"`
	IsActive *int16  `json:"is_active,omitempty"`
}

type ProductResponse struct {
	ID        int     `json:"id"`
	SKU       string  `json:"sku"`
	Name      string  `json:"name"`
	Unit      string  `json:"unit"`
	IsActive  int16   `json:"is_active"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

type ProductListResponse struct {
	Data      []ProductResponse `json:"data"`
	TotalData int               `json:"total_data"`
}
