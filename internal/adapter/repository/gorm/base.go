// internal/adapter/repository/base.go
package repository

import (
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"gorm.io/gorm"
)

type baseRepository struct {
	db *gorm.DB
}

type repositoryContainer struct {
	db                  *gorm.DB
	userRepo            repository.UserRepository
	categoryRepo        repository.CategoryRepository
	menuItemRepo        repository.MenuItemRepository
	menuOptionRepo      repository.MenuOptionRepository
	optionValueRepo     repository.OptionValueRepository
	menuItemOptionRepo  repository.MenuItemOptionRepository
	tableRepo           repository.TableRepository
	orderRepo           repository.OrderRepository
	orderItemRepo       repository.OrderItemRepository
	orderItemOptionRepo repository.OrderItemOptionRepository
	paymentRepo         repository.PaymentRepository
	revenueRepo         repository.RevenueRepository
	kitchenRepo         repository.KitchenStationRepository

	txRepo repository.TxManager
}

func NewRepositoryContainer(db *gorm.DB) repository.Repository {
	return &repositoryContainer{
		db:                  db,
		userRepo:            NewUserRepository(db),
		categoryRepo:        NewCategoryRepository(db),
		menuItemRepo:        NewMenuItemRepository(db),
		menuOptionRepo:      NewMenuOptionRepository(db),
		optionValueRepo:     NewOptionValueRepository(db),
		menuItemOptionRepo:  NewMenuItemOptionRepository(db),
		tableRepo:           NewTableRepository(db),
		orderRepo:           NewOrderRepository(db),
		orderItemRepo:       NewOrderItemRepository(db),
		orderItemOptionRepo: NewOrderItemOptionRepository(db),
		paymentRepo:         NewPaymentRepository(db),
		revenueRepo:         NewRevenueRepository(db),
		kitchenRepo:         NewKitchenStationRepository(db),
		txRepo:              NewTxManagerGorm(db),
	}
}

func (r *repositoryContainer) UserRepository() repository.UserRepository {
	return r.userRepo
}

func (r *repositoryContainer) CategoryRepository() repository.CategoryRepository {
	return r.categoryRepo
}

func (r *repositoryContainer) MenuItemRepository() repository.MenuItemRepository {
	return r.menuItemRepo
}

func (r *repositoryContainer) MenuOptionRepository() repository.MenuOptionRepository {
	return r.menuOptionRepo
}

func (r *repositoryContainer) OptionValueRepository() repository.OptionValueRepository {
	return r.optionValueRepo
}

func (r *repositoryContainer) MenuItemOptionRepository() repository.MenuItemOptionRepository {
	return r.menuItemOptionRepo
}

func (r *repositoryContainer) TableRepository() repository.TableRepository {
	return r.tableRepo
}

func (r *repositoryContainer) OrderRepository() repository.OrderRepository {
	return r.orderRepo
}

func (r *repositoryContainer) OrderItemRepository() repository.OrderItemRepository {
	return r.orderItemRepo
}

func (r *repositoryContainer) OrderItemOptionRepository() repository.OrderItemOptionRepository {
	return r.orderItemOptionRepo
}

func (r *repositoryContainer) PaymentRepository() repository.PaymentRepository {
	return r.paymentRepo
}

func (r *repositoryContainer) RevenueRepository() repository.RevenueRepository {
	return r.revenueRepo
}
func (r *repositoryContainer) KitchenStationRepository() repository.KitchenStationRepository {
	return r.kitchenRepo
}

func (r *repositoryContainer) TxManager() repository.TxManager {
	return r.txRepo
}
