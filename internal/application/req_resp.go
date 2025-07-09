package usecase

import "time"

// Request DTOs

// RegisterRequest represents user registration request
type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	Role     string `json:"role" validate:"required,oneof=candidate company_hr admin"`
}

// LoginRequest represents user login request
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// UpdateProfileRequest represents user profile update request
type UpdateProfileRequest struct {
	Email string `json:"email,omitempty" validate:"omitempty,email"`
}

// ChangePasswordRequest represents password change request
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=NewPassword"`
}

// Response DTOs

// UserResponse represents user data in responses
type UserResponse struct {
	ID            int        `json:"id"`
	Email         string     `json:"email"`
	Role          string     `json:"role"`
	IsActive      bool       `json:"is_active"`
	EmailVerified bool       `json:"email_verified"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
}

// LoginResponse represents login response with token
type LoginResponse struct {
	User  *UserResponse `json:"user"`
	Token string        `json:"token"`
}

// Category DTOs
type CreateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
}

type CategoryResponse struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

// internal/application/dto/menu_item_dto.go

// Menu Item DTOs
type CreateMenuItemRequest struct {
	CategoryID  int     `json:"category_id" validate:"required,gt=0"`
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price" validate:"required,gte=0"`
}

type UpdateMenuItemRequest struct {
	CategoryID  int     `json:"category_id" validate:"required,gt=0"`
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price" validate:"required,gte=0"`
}

type MenuItemResponse struct {
	ID          int               `json:"id"`
	CategoryID  int               `json:"category_id"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Price       float64           `json:"price"`
	CreatedAt   time.Time         `json:"created_at"`
	Category    *CategoryResponse `json:"category,omitempty"`
}

type MenuItemListResponse struct {
	Items  []*MenuItemResponse `json:"items"`
	Total  int                 `json:"total"`
	Limit  int                 `json:"limit"`
	Offset int                 `json:"offset"`
}

// internal/application/dto/table_dto.go

type CreateTableRequest struct {
	TableNumber int `json:"table_number" validate:"required,gt=0"`
	Seating     int `json:"seating" validate:"gte=0"`
}

type UpdateTableRequest struct {
	TableNumber int `json:"table_number" validate:"required,gt=0"`
	Seating     int `json:"seating" validate:"gte=0"`
}

type TableResponse struct {
	ID          int    `json:"id"`
	TableNumber int    `json:"table_number"`
	QRCode      string `json:"qr_code"`
	Seating     int    `json:"seating"`
}

// internal/application/dto/order_dto.go

// Order DTOs
type CreateOrderRequest struct {
	TableID int `json:"table_id" validate:"required,gt=0"`
}

type UpdateOrderRequest struct {
	Status string `json:"status" validate:"required,oneof=open closed"`
}

type OrderResponse struct {
	ID        int            `json:"id"`
	TableID   int            `json:"table_id"`
	Status    string         `json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	ClosedAt  *time.Time     `json:"closed_at,omitempty"`
	Table     *TableResponse `json:"table,omitempty"`
}

type OrderWithItemsResponse struct {
	ID        int                  `json:"id"`
	TableID   int                  `json:"table_id"`
	Status    string               `json:"status"`
	Items     []*OrderItemResponse `json:"items"`
	Total     float64              `json:"total"`
	CreatedAt time.Time            `json:"created_at"`
	ClosedAt  *time.Time           `json:"closed_at,omitempty"`
	Table     *TableResponse       `json:"table,omitempty"`
}

type OrderListResponse struct {
	Orders []*OrderResponse `json:"orders"`
	Total  int              `json:"total"`
	Limit  int              `json:"limit"`
	Offset int              `json:"offset"`
}

// Order Item DTOs
type AddOrderItemRequest struct {
	OrderID  int `json:"order_id" validate:"required,gt=0"`
	ItemID   int `json:"item_id" validate:"required,gt=0"`
	Quantity int `json:"quantity" validate:"required,gt=0"`
}

type UpdateOrderItemRequest struct {
	Quantity int `json:"quantity" validate:"required,gt=0"`
}

type OrderItemResponse struct {
	ID        int               `json:"id"`
	OrderID   int               `json:"order_id"`
	ItemID    int               `json:"item_id"`
	Quantity  int               `json:"quantity"`
	UnitPrice float64           `json:"unit_price"`
	Subtotal  float64           `json:"subtotal"`
	CreatedAt time.Time         `json:"created_at"`
	MenuItem  *MenuItemResponse `json:"menu_item,omitempty"`
}

type OrderTotalResponse struct {
	OrderID   int                  `json:"order_id"`
	Items     []*OrderItemResponse `json:"items"`
	Total     float64              `json:"total"`
	ItemCount int                  `json:"item_count"`
}

// internal/application/dto/payment_dto.go

// Payment DTOs
type ProcessPaymentRequest struct {
	OrderID int     `json:"order_id" validate:"required,gt=0"`
	Amount  float64 `json:"amount" validate:"required,gt=0"`
	Method  string  `json:"method" validate:"required,oneof=cash credit_card wallet"`
}

type PaymentResponse struct {
	ID      int            `json:"id"`
	OrderID int            `json:"order_id"`
	Amount  float64        `json:"amount"`
	Method  string         `json:"method"`
	PaidAt  time.Time      `json:"paid_at"`
	Order   *OrderResponse `json:"order,omitempty"`
}

type PaymentListResponse struct {
	Payments []*PaymentResponse `json:"payments"`
	Total    int                `json:"total"`
	Limit    int                `json:"limit"`
	Offset   int                `json:"offset"`
}

// internal/application/dto/revenue_dto.go

// Revenue DTOs
type DailyRevenueResponse struct {
	Date         time.Time `json:"date"`
	TotalRevenue float64   `json:"total_revenue"`
	OrderCount   int       `json:"order_count,omitempty"`
}

type MonthlyRevenueResponse struct {
	Month        time.Time `json:"month"`
	TotalRevenue float64   `json:"total_revenue"`
	OrderCount   int       `json:"order_count,omitempty"`
}

type TotalRevenueResponse struct {
	StartDate    time.Time `json:"start_date"`
	EndDate      time.Time `json:"end_date"`
	TotalRevenue float64   `json:"total_revenue"`
	OrderCount   int       `json:"order_count,omitempty"`
}

// internal/application/dto/qr_code_dto.go

type QRCodeScanResponse struct {
	TableID      int            `json:"table_id"`
	Table        *TableResponse `json:"table"`
	HasOpenOrder bool           `json:"has_open_order"`
	OpenOrder    *OrderResponse `json:"open_order,omitempty"`
}

// internal/application/dto/common_dto.go

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

type SuccessResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

// Pagination
type PaginationRequest struct {
	Limit  int `json:"limit" validate:"min=1,max=100"`
	Offset int `json:"offset" validate:"min=0"`
}

type PaginationResponse struct {
	Total   int  `json:"total"`
	Limit   int  `json:"limit"`
	Offset  int  `json:"offset"`
	HasNext bool `json:"has_next"`
}
