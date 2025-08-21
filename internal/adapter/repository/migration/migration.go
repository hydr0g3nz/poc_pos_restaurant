// internal/adapter/repository/migration/migration.go
package migration

import (
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/repository/gorm/model"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&model.User{},
		&model.Category{},
		&model.MenuItem{},
		&model.MenuOption{},
		&model.OptionValue{},
		&model.MenuItemOption{},
		&model.Table{},
		&model.Order{},
		&model.OrderItem{},
		&model.OrderItemOption{},
		&model.Payment{},
	)
}
