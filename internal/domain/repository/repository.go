// internal/domain/repository/repository.go (updated)
package repository

import (
	"context"
	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
)

type Repository interface {
	UserRepository() UserRepository
	CategoryRepository() CategoryRepository
	MenuItemRepository() MenuItemRepository
	TableRepository() TableRepository
	OrderRepository() OrderRepository
	OrderItemRepository() OrderItemRepository
	PaymentRepository() PaymentRepository
	RevenueRepository() RevenueRepository
}

type DBTransaction interface {
	DoInTransaction(fn func(repo Repository) error) error
}

// Existing UserRepository
type UserRepository interface {
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	GetByID(ctx context.Context, id int) (*entity.User, error)
	GetByEmail(ctx context.Context, email string) (*entity.User, error)
	Update(ctx context.Context, user *entity.User) (*entity.User, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*entity.User, error)
	ListByRole(ctx context.Context, role string, limit, offset int) ([]*entity.User, error)
	UpdateLastLogin(ctx context.Context, id int) error
}

// CategoryRepository handles category operations
type CategoryRepository interface {
	Create(ctx context.Context, category *entity.Category) (*entity.Category, error)
	GetByID(ctx context.Context, id int) (*entity.Category, error)
	GetByName(ctx context.Context, name string) (*entity.Category, error)
	Update(ctx context.Context, category *entity.Category) (*entity.Category, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]*entity.Category, error)
	HasMenuItems(ctx context.Context, categoryID int) (bool, error)
}

// MenuItemRepository handles menu item operations
type MenuItemRepository interface {
	Create(ctx context.Context, item *entity.MenuItem) (*entity.MenuItem, error)
	GetByID(ctx context.Context, id int) (*entity.MenuItem, error)
	Update(ctx context.Context, item *entity.MenuItem) (*entity.MenuItem, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*entity.MenuItem, error)
	ListByCategory(ctx context.Context, categoryID int, limit, offset int) ([]*entity.MenuItem, error)
	Search(ctx context.Context, query string, limit, offset int) ([]*entity.MenuItem, error)
}

// TableRepository handles table operations
type TableRepository interface {
	Create(ctx context.Context, table *entity.Table) (*entity.Table, error)
	GetByID(ctx context.Context, id int) (*entity.Table, error)
	GetByNumber(ctx context.Context, number int) (*entity.Table, error)
	GetByQRCode(ctx context.Context, qrCode string) (*entity.Table, error)
	Update(ctx context.Context, table *entity.Table) (*entity.Table, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context) ([]*entity.Table, error)
	HasOrders(ctx context.Context, tableID int) (bool, error)
}

// OrderRepository handles order operations
type OrderRepository interface {
	Create(ctx context.Context, order *entity.Order) (*entity.Order, error)
	GetByID(ctx context.Context, id int) (*entity.Order, error)
	GetByIDWithItems(ctx context.Context, id int) (*entity.Order, error)
	Update(ctx context.Context, order *entity.Order) (*entity.Order, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*entity.Order, error)
	ListByTable(ctx context.Context, tableID int, limit, offset int) ([]*entity.Order, error)
	GetOpenOrderByTable(ctx context.Context, tableID int) (*entity.Order, error)
	ListByStatus(ctx context.Context, status string, limit, offset int) ([]*entity.Order, error)
	ListByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*entity.Order, error)
}

// OrderItemRepository handles order item operations
type OrderItemRepository interface {
	Create(ctx context.Context, item *entity.OrderItem) (*entity.OrderItem, error)
	GetByID(ctx context.Context, id int) (*entity.OrderItem, error)
	Update(ctx context.Context, item *entity.OrderItem) (*entity.OrderItem, error)
	Delete(ctx context.Context, id int) error
	ListByOrder(ctx context.Context, orderID int) ([]*entity.OrderItem, error)
	DeleteByOrder(ctx context.Context, orderID int) error
	GetByOrderAndItem(ctx context.Context, orderID, itemID int) (*entity.OrderItem, error)
}

// PaymentRepository handles payment operations
type PaymentRepository interface {
	Create(ctx context.Context, payment *entity.Payment) (*entity.Payment, error)
	GetByID(ctx context.Context, id int) (*entity.Payment, error)
	GetByOrderID(ctx context.Context, orderID int) (*entity.Payment, error)
	Update(ctx context.Context, payment *entity.Payment) (*entity.Payment, error)
	Delete(ctx context.Context, id int) error
	List(ctx context.Context, limit, offset int) ([]*entity.Payment, error)
	ListByDateRange(ctx context.Context, startDate, endDate time.Time, limit, offset int) ([]*entity.Payment, error)
	ListByMethod(ctx context.Context, method string, limit, offset int) ([]*entity.Payment, error)
}

// RevenueRepository handles revenue reporting
type RevenueRepository interface {
	GetDailyRevenue(ctx context.Context, date time.Time) (*entity.DailyRevenue, error)
	GetMonthlyRevenue(ctx context.Context, year int, month int) (*entity.MonthlyRevenue, error)
	GetDailyRevenueRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.DailyRevenue, error)
	GetMonthlyRevenueRange(ctx context.Context, startDate, endDate time.Time) ([]*entity.MonthlyRevenue, error)
	GetTotalRevenue(ctx context.Context, startDate, endDate time.Time) (float64, error)
}
