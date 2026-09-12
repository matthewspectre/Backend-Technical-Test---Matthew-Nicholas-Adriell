package warehouse

type WarehouseCreateRequest struct {
	Code     string `json:"code"`
	Name     string `json:"name"`
	Location string `json:"location"`
	IsActive int16  `json:"is_active"`
}

type WarehouseUpdateRequest struct {
	Code     *string `json:"code,omitempty"`
	Name     *string `json:"name,omitempty"`
	Location *string `json:"location,omitempty"`
	IsActive *int16  `json:"is_active,omitempty"`
}

type WarehouseResponse struct {
	ID        int     `json:"id"`
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Location  string  `json:"location"`
	IsActive  int16   `json:"is_active"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt *string `json:"updated_at"`
}

type WarehouseListResponse struct {
	Data      []WarehouseResponse `json:"data"`
	TotalData int                 `json:"total_data"`
}
