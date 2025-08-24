package usecase

import (
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
)

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
	Name         string `json:"name" validate:"required,min=1,max=50"`
	Description  string `json:"description,omitempty" validate:"max=100"`
	DisplayOrder int    `json:"display_order,omitempty"`
	IsActive     bool   `json:"is_active"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" validate:"required,min=1,max=50"`
}

type CategoryResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	DisplayOrder int    `json:"display_order"`
	IsActive     bool   `json:"is_active"`
	// CreatedAt time.Time `json:"created_at"`

}

// internal/application/dto/menu_item_dto.go

// Menu Item DTOs
type CreateMenuItemRequest struct {
	CategoryID       int     `json:"category_id" validate:"required,gt=0"`
	KitchenStationID int     `json:"kitchen_station_id" validate:"required,gt=0"`
	Name             string  `json:"name" validate:"required,min=1,max=100"`
	Description      string  `json:"description,omitempty"`
	Price            float64 `json:"price" validate:"required,gte=0"`
}

type UpdateMenuItemRequest struct {
	CategoryID  int     `json:"category_id" validate:"required,gt=0"`
	Name        string  `json:"name" validate:"required,min=1,max=100"`
	Description string  `json:"description,omitempty"`
	Price       float64 `json:"price" validate:"required,gte=0"`
}

type MenuItemResponse struct {
	ID          int     `json:"id"`
	CategoryID  int     `json:"category_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	// CreatedAt      time.Time                   `json:"created_at"`
	Category       string                   `json:"category"`
	KitchenStation string                   `json:"kitchen_station"`
	IsActive       bool                     `json:"is_active"`
	IsRecommended  bool                     `json:"is_recommended"`
	DisplayOrder   int                      `json:"display_order"`
	MenuOption     []*entity.MenuItemOption `json:"menu_option"`
	// DiscountPercent float64 `json:"discount_percent"`
	// IsDiscounted    bool    `json:"is_discounted"`
	// Category       *CategoryResponse           `json:"category,omitempty"`
	// KitchenStation *KitchenStationOnlyResponse `json:"kitchen_station,omitempty"`
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
type OrderWithItemsListResponse struct {
	Orders []*OrderWithItemsResponse `json:"orders"`
	Total  int                       `json:"total"`
	Limit  int                       `json:"limit"`
	Offset int                       `json:"offset"`
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
	ID             int               `json:"id"`
	OrderID        int               `json:"order_id"`
	ItemID         int               `json:"item_id"`
	Quantity       int               `json:"quantity"`
	UnitPrice      float64           `json:"unit_price"`
	Subtotal       float64           `json:"subtotal"`
	CreatedAt      time.Time         `json:"created_at"`
	MenuItem       *MenuItemResponse `json:"menu_item,omitempty"`
	Name           string            `json:"name"`
	KitchenStation string            `json:"kitchen_station,omitempty"` // optional kitchen ID for tracking

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

// Menu Option DTOs
type CreateMenuOptionRequest struct {
	Name       string `json:"name" validate:"required,min=1,max=100"`
	Type       string `json:"type" validate:"required,min=1,max=50"` // เช่น "spice_level", "temperature", "size"
	IsRequired bool   `json:"is_required"`
}

type UpdateMenuOptionRequest struct {
	Name       string `json:"name" validate:"required,min=1,max=100"`
	Type       string `json:"type" validate:"required,min=1,max=50"`
	IsRequired bool   `json:"is_required"`
}

type MenuOptionResponse struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Type       string `json:"type"`
	IsRequired bool   `json:"is_required"`
}

type MenuOptionListResponse struct {
	Options []*MenuOptionResponse `json:"options"`
	Total   int                   `json:"total"`
	Limit   int                   `json:"limit"`
	Offset  int                   `json:"offset"`
}

// Option Value DTOs
type CreateOptionValueRequest struct {
	OptionID        int     `json:"option_id" validate:"required,gt=0"`
	Name            string  `json:"name" validate:"required,min=1,max=100"`
	IsDefault       bool    `json:"is_default"`
	AdditionalPrice float64 `json:"additional_price,omitempty" validate:"gte=0"`
	DisplayOrder    int     `json:"display_order,omitempty"`
}

type UpdateOptionValueRequest struct {
	Name            string  `json:"name" validate:"required,min=1,max=100"`
	IsDefault       bool    `json:"is_default"`
	AdditionalPrice float64 `json:"additional_price,omitempty" validate:"gte=0"`
	DisplayOrder    int     `json:"display_order,omitempty"`
}

type OptionValueResponse struct {
	ID              int                 `json:"id"`
	OptionID        int                 `json:"option_id"`
	Name            string              `json:"name"`
	IsDefault       bool                `json:"is_default"`
	AdditionalPrice float64             `json:"additional_price"`
	DisplayOrder    int                 `json:"display_order"`
	Option          *MenuOptionResponse `json:"option,omitempty"`
}

type OptionValueListResponse struct {
	Values []*OptionValueResponse `json:"values"`
	Total  int                    `json:"total"`
	Limit  int                    `json:"limit"`
	Offset int                    `json:"offset"`
}

// Menu Item Option DTOs
type AddMenuItemOptionRequest struct {
	ItemID   int  `json:"item_id" validate:"required,gt=0"`
	OptionID int  `json:"option_id" validate:"required,gt=0"`
	IsActive bool `json:"is_active"`
}

type UpdateMenuItemOptionRequest struct {
	IsActive bool `json:"is_active"`
}

type MenuItemOptionResponse struct {
	ItemID   int                    `json:"item_id"`
	OptionID int                    `json:"option_id"`
	IsActive bool                   `json:"is_active"`
	Option   *MenuOptionResponse    `json:"option,omitempty"`
	Values   []*OptionValueResponse `json:"values,omitempty"`
}

type MenuItemOptionListResponse struct {
	ItemOptions []*MenuItemOptionResponse `json:"item_options"`
	Total       int                       `json:"total"`
}

// Order Item Option DTOs
type AddOrderItemOptionRequest struct {
	OrderItemID     int     `json:"order_item_id" validate:"required,gt=0"`
	OptionID        int     `json:"option_id" validate:"required,gt=0"`
	ValueID         int     `json:"value_id" validate:"required,gt=0"`
	AdditionalPrice float64 `json:"additional_price,omitempty" validate:"gte=0"`
}

type UpdateOrderItemOptionRequest struct {
	ValueID         int     `json:"value_id" validate:"required,gt=0"`
	AdditionalPrice float64 `json:"additional_price,omitempty" validate:"gte=0"`
}

type OrderItemOptionResponse struct {
	OrderItemID     int                  `json:"order_item_id"`
	OptionID        int                  `json:"option_id"`
	ValueID         int                  `json:"value_id"`
	AdditionalPrice float64              `json:"additional_price"`
	Option          *MenuOptionResponse  `json:"option,omitempty"`
	Value           *OptionValueResponse `json:"value,omitempty"`
}

type OrderItemOptionListResponse struct {
	ItemOptions []*OrderItemOptionResponse `json:"item_options"`
	Total       int                        `json:"total"`
}

// Kitchen Management DTOs
type KitchenOrderResponse struct {
	OrderID       int                         `json:"order_id"`
	OrderNumber   int                         `json:"order_number"`
	TableNumber   *int                        `json:"table_number,omitempty"`
	CustomerName  string                      `json:"customer_name,omitempty"`
	OrderType     string                      `json:"order_type"`
	Items         []*KitchenOrderItemResponse `json:"items"`
	CreatedAt     time.Time                   `json:"created_at"`
	EstimatedTime int                         `json:"estimated_time,omitempty"` // minutes
}

type KitchenOrderItemResponse struct {
	ID              int                        `json:"id"`
	ItemID          int                        `json:"item_id"`
	Name            string                     `json:"name"`
	Quantity        int                        `json:"quantity"`
	Status          string                     `json:"status"`
	PreparationTime int                        `json:"preparation_time,omitempty"`
	KitchenStation  string                     `json:"kitchen_station,omitempty"`
	KitchenNotes    string                     `json:"kitchen_notes,omitempty"`
	Notes           string                     `json:"notes,omitempty"`
	Options         []*OrderItemOptionResponse `json:"options,omitempty"`
	CreatedAt       time.Time                  `json:"created_at"`
	StartedAt       *time.Time                 `json:"started_at,omitempty"`
	ReadyAt         *time.Time                 `json:"ready_at,omitempty"`
	ServedAt        *time.Time                 `json:"served_at,omitempty"`
}

type UpdateOrderItemStatusRequest struct {
	Status       string `json:"status" validate:"required,oneof=pending preparing ready served cancelled"`
	KitchenNotes string `json:"kitchen_notes,omitempty" validate:"max=200"`
}

type KitchenQueueResponse struct {
	Queue          []*KitchenOrderResponse `json:"queue"`
	TotalItems     int                     `json:"total_items"`
	PendingItems   int                     `json:"pending_items"`
	PreparingItems int                     `json:"preparing_items"`
	ReadyItems     int                     `json:"ready_items"`
}

type KitchenStationResponse struct {
	Station     string                  `json:"station"`
	Orders      []*KitchenOrderResponse `json:"orders"`
	TotalItems  int                     `json:"total_items"`
	AverageTime int                     `json:"average_preparation_time"` // minutes
}

// Enhanced Search/Filter DTOs
type OrderFilterRequest struct {
	PaginationRequest
	Status        string     `json:"status,omitempty" validate:"omitempty,oneof=open ordered completed cancelled"`
	OrderType     string     `json:"order_type,omitempty" validate:"omitempty,oneof=dine_in phone online"`
	PaymentStatus string     `json:"payment_status,omitempty" validate:"omitempty,oneof=unpaid paid refunded"`
	TableID       *int       `json:"table_id,omitempty" validate:"omitempty,gt=0"`
	StartDate     *time.Time `json:"start_date,omitempty"`
	EndDate       *time.Time `json:"end_date,omitempty"`
	CustomerName  string     `json:"customer_name,omitempty" validate:"max=100"`
	CustomerPhone string     `json:"customer_phone,omitempty" validate:"max=20"`
}

type MenuItemFilterRequest struct {
	PaginationRequest
	CategoryID    *int     `json:"category_id,omitempty" validate:"omitempty,gt=0"`
	IsActive      *bool    `json:"is_active,omitempty"`
	MinPrice      *float64 `json:"min_price,omitempty" validate:"omitempty,gte=0"`
	MaxPrice      *float64 `json:"max_price,omitempty" validate:"omitempty,gte=0"`
	Search        string   `json:"search,omitempty" validate:"max=100"`
	IsRecommended *bool    `json:"is_recommended,omitempty"`
}

type CreateKitchenStationRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	IsAvailable bool   `json:"is_available" validate:"required"`
}

type UpdateKitchenStationRequest struct {
	Name        string `json:"name" validate:"required,min=1,max=100"`
	IsAvailable bool   `json:"is_available" validate:"required"`
}
type KitchenStationOnlyResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	IsAvailable bool   `json:"is_available"`
}

type AddOrderItemListRequest struct {
	OrderID int                 `json:"order_id" validate:"required,gt=0"`
	Items   []*OrderItemRequest `json:"items" validate:"required,dive,required"`
}
type OrderItemRequest struct {
	MenuItemID int                       `json:"menu_item_id" validate:"required,gt=0"`
	Quantity   int                       `json:"quantity" validate:"required,gt=0"`
	Options    []*OrderItemOptionRequest `json:"options,omitempty"`
}

type OrderItemOptionRequest struct {
	OptionID    int `json:"option_id" validate:"required,gt=0"`
	OptionValID int `json:"option_val_id" validate:"required,gt=0"`
}

type UpdateOrderItemListRequest struct {
	OrderID int                        `json:"order_id" validate:"required,gt=0"`
	Items   []*UpdateOrderItemRequest2 `json:"items" validate:"required,dive,required"`
}

// เพิ่มใน internal/application/req_resp.go

type UpdateOrderItemRequest2 struct {
	OrderItemID int                             `json:"order_item_id" validate:"required,gt=0"`
	MenuItemID  int                             `json:"menu_item_id" validate:"required,gt=0"`
	Quantity    int                             `json:"quantity" validate:"required,gt=0"`
	Options     []*OrderItemOptionUpdateRequest `json:"options,omitempty"`
	Action      string                          `json:"action,omitempty" validate:"omitempty,oneof=update delete"`
}

type OrderItemOptionUpdateRequest struct {
	OptionID    int    `json:"option_id" validate:"required,gt=0"`
	OptionValID int    `json:"option_val_id" validate:"required,gt=0"`
	Action      string `json:"action,omitempty" validate:"omitempty,oneof=add update delete"` // add, update, delete
}
