// internal/adapter/repository/model/models.go
package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID            int       `gorm:"primaryKey;autoIncrement"`
	Email         string    `gorm:"uniqueIndex;not null"`
	PasswordHash  string    `gorm:"not null"`
	Role          string    `gorm:"not null"`
	IsActive      bool      `gorm:"default:true"`
	EmailVerified bool      `gorm:"default:false"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`
	LastLoginAt   *time.Time
	DeletedAt     gorm.DeletedAt `gorm:"index"`
}

type Category struct {
	ID           int    `gorm:"primaryKey;autoIncrement"`
	Name         string `gorm:"uniqueIndex;not null"`
	Description  string
	DisplayOrder int            `gorm:"default:0"`
	IsActive     bool           `gorm:"default:true"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`

	// Relationships
	MenuItems []MenuItem `gorm:"foreignKey:CategoryID"`
}

type MenuItem struct {
	ID              int    `gorm:"primaryKey;autoIncrement"`
	CategoryID      int    `gorm:"not null;index"`
	Name            string `gorm:"not null"`
	Description     string
	Price           int64 `gorm:"not null"` // stored in satang
	ImageURL        string
	IsRecommended   bool `gorm:"default:false"`
	PreparationTime int  `gorm:"default:0"` // in minutes
	DisplayOrder    int  `gorm:"default:0"`
	KitchenID       string
	IsActive        bool           `gorm:"default:true"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	// Relationships
	Category        Category         `gorm:"foreignKey:CategoryID"`
	OrderItems      []OrderItem      `gorm:"foreignKey:ItemID"`
	MenuItemOptions []MenuItemOption `gorm:"foreignKey:ItemID"`
}

type MenuOption struct {
	ID         int            `gorm:"primaryKey;autoIncrement"`
	Name       string         `gorm:"not null"`
	Type       string         `gorm:"not null"`
	IsRequired bool           `gorm:"default:false"`
	CreatedAt  time.Time      `gorm:"autoCreateTime"`
	UpdatedAt  time.Time      `gorm:"autoUpdateTime"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`

	// Relationships
	OptionValues     []OptionValue     `gorm:"foreignKey:OptionID"`
	MenuItemOptions  []MenuItemOption  `gorm:"foreignKey:OptionID"`
	OrderItemOptions []OrderItemOption `gorm:"foreignKey:OptionID"`
}

type OptionValue struct {
	ID              int            `gorm:"primaryKey;autoIncrement"`
	OptionID        int            `gorm:"not null;index"`
	Name            string         `gorm:"not null"`
	IsDefault       bool           `gorm:"default:false"`
	AdditionalPrice int64          `gorm:"default:0"` // stored in satang
	DisplayOrder    int            `gorm:"default:0"`
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	// Relationships
	MenuOption       MenuOption        `gorm:"foreignKey:OptionID"`
	OrderItemOptions []OrderItemOption `gorm:"foreignKey:ValueID"`
}

type MenuItemOption struct {
	ItemID   int  `gorm:"primaryKey"`
	OptionID int  `gorm:"primaryKey"`
	IsActive bool `gorm:"default:true"`

	// Relationships
	MenuItem   MenuItem   `gorm:"foreignKey:ItemID"`
	MenuOption MenuOption `gorm:"foreignKey:OptionID"`
}

type Table struct {
	ID          int            `gorm:"primaryKey;autoIncrement"`
	TableNumber int            `gorm:"uniqueIndex;not null"`
	Seating     int            `gorm:"not null"`
	IsActive    bool           `gorm:"default:true"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	// Relationships
	Orders []Order `gorm:"foreignKey:TableID"`
}

type Order struct {
	ID                  int    `gorm:"primaryKey;autoIncrement"`
	OrderNumber         int    `gorm:"uniqueIndex;autoIncrement"`
	TableID             int    `gorm:"not null;index"`
	OrderStatus         string `gorm:"not null;default:'open'"`
	PaymentStatus       string `gorm:"not null;default:'unpaid'"`
	QRCode              string `gorm:"uniqueIndex"`
	Notes               string
	SpecialInstructions string
	Subtotal            int64     `gorm:"default:0"` // stored in satang
	Discount            int64     `gorm:"default:0"` // stored in satang
	TaxAmount           int64     `gorm:"default:0"` // stored in satang
	ServiceCharge       int64     `gorm:"default:0"` // stored in satang
	Total               int64     `gorm:"default:0"` // stored in satang
	CreatedAt           time.Time `gorm:"autoCreateTime"`
	UpdatedAt           time.Time `gorm:"autoUpdateTime"`
	ClosedAt            *time.Time
	DeletedAt           gorm.DeletedAt `gorm:"index"`

	// Relationships
	Table      Table       `gorm:"foreignKey:TableID"`
	OrderItems []OrderItem `gorm:"foreignKey:OrderID"`
	Payments   []Payment   `gorm:"foreignKey:OrderID"`
}

type OrderItem struct {
	ID              int    `gorm:"primaryKey;autoIncrement"`
	OrderID         int    `gorm:"not null;index"`
	ItemID          int    `gorm:"not null;index"`
	Quantity        int    `gorm:"not null"`
	UnitPrice       int64  `gorm:"not null"` // stored in satang
	Name            string `gorm:"not null"`
	Discount        int64  `gorm:"default:0"` // stored in satang
	Total           int64  `gorm:"not null"`  // stored in satang
	SpecialReq      string
	ItemStatus      string `gorm:"not null;default:'pending'"`
	OrderNumber     string
	KitchenTicketID int
	KitchenStation  string
	KitchenNotes    string
	ServedAt        *time.Time
	CreatedAt       time.Time      `gorm:"autoCreateTime"`
	UpdatedAt       time.Time      `gorm:"autoUpdateTime"`
	DeletedAt       gorm.DeletedAt `gorm:"index"`

	// Relationships
	Order            Order             `gorm:"foreignKey:OrderID"`
	MenuItem         MenuItem          `gorm:"foreignKey:ItemID"`
	OrderItemOptions []OrderItemOption `gorm:"foreignKey:OrderItemID"`
}

type OrderItemOption struct {
	OrderItemID     int   `gorm:"primaryKey"`
	OptionID        int   `gorm:"primaryKey"`
	ValueID         int   `gorm:"primaryKey"`
	AdditionalPrice int64 `gorm:"default:0"` // stored in satang

	// Relationships
	OrderItem   OrderItem   `gorm:"foreignKey:OrderItemID"`
	MenuOption  MenuOption  `gorm:"foreignKey:OptionID"`
	OptionValue OptionValue `gorm:"foreignKey:ValueID"`
}

type Payment struct {
	ID        int    `gorm:"primaryKey;autoIncrement"`
	OrderID   int    `gorm:"not null;index"`
	Amount    int64  `gorm:"not null"` // stored in satang
	Method    string `gorm:"not null"`
	Reference string
	PaidAt    time.Time      `gorm:"autoCreateTime"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	// Relationships
	Order Order `gorm:"foreignKey:OrderID"`
}
