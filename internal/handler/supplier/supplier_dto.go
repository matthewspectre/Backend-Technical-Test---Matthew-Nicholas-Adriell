package supplier

type SupplierCreateRequest struct {
	CompanyName string `json:"company_name"`
	Name        string `json:"name"`
	Email       string `json:"email"`
	Phone       string `json:"phone"`
	Address     string `json:"address"`
	IsActive    int16  `json:"is_active"`
}

type SupplierUpdateRequest struct {
	CompanyName *string `json:"company_name,omitempty"`
	Name        *string `json:"name,omitempty"`
	Email       *string `json:"email,omitempty"`
	Phone       *string `json:"phone,omitempty"`
	Address     *string `json:"address,omitempty"`
	IsActive    *int16  `json:"is_active,omitempty"`
}

type SupplierResponse struct {
	ID          int     `json:"id"`
	CompanyName string  `json:"company_name"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	Phone       string  `json:"phone"`
	Address     string  `json:"address"`
	IsActive    int16   `json:"is_active"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   *string `json:"updated_at"`
}

type SupplierListResponse struct {
	Data      []SupplierResponse `json:"data"`
	TotalData int                `json:"total_data"`
}
