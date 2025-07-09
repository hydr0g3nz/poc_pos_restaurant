package repository

import (
	"context"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
)

type Repository interface {
	UserRepository() UserRepository
}
type DBTransaction interface {
	DoInTransaction(fn func(repo Repository) error) error
}

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
