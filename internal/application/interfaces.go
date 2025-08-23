// internal/application/interfaces.go
package usecase

import (
	"context"
	"time"
)

// CategoryUsecase handles category business logic
type CategoryUsecase interface {
	CreateCategory(ctx context.Context, req *CreateCategoryRequest) (*CategoryResponse, error)
	GetCategory(ctx context.Context, id int) (*CategoryResponse, error)
	GetCategoryByName(ctx context.Context, name string) (*CategoryResponse, error)
	UpdateCategory(ctx context.Context, id int, req *UpdateCategoryRequest) (*CategoryResponse, error)
	DeleteCategory(ctx context.Context, id int) error
	ListCategories(ctx context.Context) ([]*CategoryResponse, error)
}

// MenuItemUsecase handles menu item business logic
type MenuItemUsecase interface {
	CreateMenuItem(ctx context.Context, req *CreateMenuItemRequest) (*MenuItemResponse, error)
	GetMenuItem(ctx context.Context, id int) (*MenuItemResponse, error)
	UpdateMenuItem(ctx context.Context, id int, req *UpdateMenuItemRequest) (*MenuItemResponse, error)
	DeleteMenuItem(ctx context.Context, id int) error
	ListMenuItems(ctx context.Context, limit, offset int) (*MenuItemListResponse, error)
	ListMenuItemsByCategory(ctx context.Context, categoryID int, limit, offset int) (*MenuItemListResponse, error)
	SearchMenuItems(ctx context.Context, query string, limit, offset int) (*MenuItemListResponse, error)
}

// TableUsecase handles table business logic
type TableUsecase interface {
	CreateTable(ctx context.Context, req *CreateTableRequest) (*TableResponse, error)
	GetTable(ctx context.Context, id int) (*TableResponse, error)
	GetTableByNumber(ctx context.Context, number int) (*TableResponse, error)
	GetTableByQRCode(ctx context.Context, qrCode string) (*TableResponse, error)
	UpdateTable(ctx context.Context, id int, req *UpdateTableRequest) (*TableResponse, error)
	DeleteTable(ctx context.Context, id int) error
	ListTables(ctx context.Context) ([]*TableResponse, error)
}

// OrderUsecase handles order business logic
type OrderUsecase interface {
	CreateOrder(ctx context.Context, req *CreateOrderRequest) (*OrderResponse, string, error)
	GetOrder(ctx context.Context, id int) (*OrderResponse, error)
	GetOrderWithItems(ctx context.Context, id int) (*OrderWithItemsResponse, error)
	UpdateOrder(ctx context.Context, id int, req *UpdateOrderRequest) (*OrderResponse, error)
	CloseOrder(ctx context.Context, id int) (*OrderResponse, error)
	ListOrders(ctx context.Context, limit, offset int) (*OrderListResponse, error)
	ListOrdersByTable(ctx context.Context, tableID int, limit, offset int) (*OrderListResponse, error)
	GetOpenOrderByTable(ctx context.Context, tableID int) (*OrderResponse, error)
	GetOrdersByStatus(ctx context.Context, status string, limit, offset int) (*OrderListResponse, error)
	GetOrdersByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) (*OrderListResponse, error)
	PrintOrderReceipt(ctx context.Context, orderID int) error
	PrintOrderQRCode(ctx context.Context, orderID int) error
	// Order item operations
	AddOrderItem(ctx context.Context, req *AddOrderItemRequest) (*OrderItemResponse, error)
	UpdateOrderItem(ctx context.Context, id int, req *UpdateOrderItemRequest) (*OrderItemResponse, error)
	RemoveOrderItem(ctx context.Context, id int) error
	ListOrderItems(ctx context.Context, orderID int) ([]*OrderItemResponse, error)

	// Calculate bill
	CalculateOrderTotal(ctx context.Context, orderID int) (*OrderTotalResponse, error)
}

// PaymentUsecase handles payment business logic
type PaymentUsecase interface {
	ProcessPayment(ctx context.Context, req *ProcessPaymentRequest) (*PaymentResponse, error)
	GetPayment(ctx context.Context, id int) (*PaymentResponse, error)
	GetPaymentByOrder(ctx context.Context, orderID int) (*PaymentResponse, error)
	ListPayments(ctx context.Context, limit, offset int) (*PaymentListResponse, error)
	ListPaymentsByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) (*PaymentListResponse, error)
	ListPaymentsByMethod(ctx context.Context, method string, limit, offset int) (*PaymentListResponse, error)
}

// RevenueUsecase handles revenue reporting business logic
type RevenueUsecase interface {
	GetDailyRevenue(ctx context.Context, date time.Time) (*DailyRevenueResponse, error)
	GetMonthlyRevenue(ctx context.Context, year int, month int) (*MonthlyRevenueResponse, error)
	GetDailyRevenueRange(ctx context.Context, startDate, endDate time.Time) ([]*DailyRevenueResponse, error)
	GetMonthlyRevenueRange(ctx context.Context, startDate, endDate time.Time) ([]*MonthlyRevenueResponse, error)
	GetTotalRevenue(ctx context.Context, startDate, endDate time.Time) (*TotalRevenueResponse, error)
}

// QRCodeUsecase handles QR code scanning and order creation
type QRCodeUsecase interface {
	ScanQRCode(ctx context.Context, qrCode string) (*QRCodeScanResponse, error)
	CreateOrderFromQRCode(ctx context.Context, qrCode string) (*OrderResponse, error)
}

// =========================== new
// MenuOptionUsecase - จัดการตัวเลือกเมนู (เผ็ด, หวาน, เพิ่มผัก)
type MenuOptionUsecase interface {
	CreateMenuOption(ctx context.Context, req *CreateMenuOptionRequest) (*MenuOptionResponse, error)
	GetMenuOption(ctx context.Context, id int) (*MenuOptionResponse, error)
	UpdateMenuOption(ctx context.Context, id int, req *UpdateMenuOptionRequest) (*MenuOptionResponse, error)
	DeleteMenuOption(ctx context.Context, id int) error
	ListMenuOptions(ctx context.Context) ([]*MenuOptionResponse, error)
	GetMenuOptionsByType(ctx context.Context, optionType string) ([]*MenuOptionResponse, error)
}

// OptionValueUsecase - จัดการค่าตัวเลือก (เผ็ดน้อย, เผ็ดปานกลาง, เผ็ดมาก)
type OptionValueUsecase interface {
	CreateOptionValue(ctx context.Context, req *CreateOptionValueRequest) (*OptionValueResponse, error)
	GetOptionValue(ctx context.Context, id int) (*OptionValueResponse, error)
	GetOptionValuesByOptionID(ctx context.Context, optionID int) ([]*OptionValueResponse, error)
	UpdateOptionValue(ctx context.Context, id int, req *UpdateOptionValueRequest) (*OptionValueResponse, error)
	DeleteOptionValue(ctx context.Context, id int) error
}

// MenuItemOptionUsecase - จัดการตัวเลือกของเมนูแต่ละรายการ
type MenuItemOptionUsecase interface {
	AddOptionToMenuItem(ctx context.Context, req *AddMenuItemOptionRequest) (*MenuItemOptionResponse, error)
	RemoveOptionFromMenuItem(ctx context.Context, itemID, optionID int) error
	GetMenuItemOptions(ctx context.Context, itemID int) ([]*MenuItemOptionResponse, error)
	UpdateMenuItemOption(ctx context.Context, req *UpdateMenuItemOptionRequest) (*MenuItemOptionResponse, error)
}

// OrderItemOptionUsecase - จัดการตัวเลือกของรายการในคำสั่งซื้อ
type OrderItemOptionUsecase interface {
	AddOptionToOrderItem(ctx context.Context, req *AddOrderItemOptionRequest) (*OrderItemOptionResponse, error)
	UpdateOrderItemOption(ctx context.Context, req *UpdateOrderItemOptionRequest) (*OrderItemOptionResponse, error)
	RemoveOptionFromOrderItem(ctx context.Context, orderItemID, optionID, valueID int) error
	GetOrderItemOptions(ctx context.Context, orderItemID int) ([]*OrderItemOptionResponse, error)
}

// KitchenUsecase - จัดการครัว/การเตรียมอาหาร
type KitchenUsecase interface {
	GetKitchenQueue(ctx context.Context) ([]*KitchenOrderResponse, error)
	UpdateOrderItemStatus(ctx context.Context, orderItemID int, status string) (*OrderItemResponse, error)
	GetOrderItemsByStatus(ctx context.Context, status string) ([]*OrderItemResponse, error)
	MarkOrderItemAsReady(ctx context.Context, orderItemID int) (*OrderItemResponse, error)
	MarkOrderItemAsServed(ctx context.Context, orderItemID int) (*OrderItemResponse, error)
	GetKitchenOrdersByStation(ctx context.Context, station string) ([]*OrderItemResponse, error)
}
