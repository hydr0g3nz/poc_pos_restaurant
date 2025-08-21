// internal/adapter/repository/user_repository.go
package repository

import (
	"context"

	"time"

	"github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/repository/gorm/model"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
	"gorm.io/gorm"
)

type userRepository struct {
	baseRepository
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepository{
		baseRepository: baseRepository{db: db},
	}
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	dbUser := r.entityToModel(user)

	if err := r.db.WithContext(ctx).Create(dbUser).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbUser)
}

func (r *userRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
	var dbUser model.User

	if err := r.db.WithContext(ctx).First(&dbUser, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbUser)
}

func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
	var dbUser model.User

	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&dbUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}

	return r.modelToEntity(&dbUser)
}

func (r *userRepository) Update(ctx context.Context, user *entity.User) (*entity.User, error) {
	dbUser := r.entityToModel(user)

	if err := r.db.WithContext(ctx).Save(dbUser).Error; err != nil {
		return nil, err
	}

	return r.modelToEntity(dbUser)
}

func (r *userRepository) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&model.User{}, id).Error
}

func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entity.User, error) {
	var dbUsers []model.User

	query := r.db.WithContext(ctx)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&dbUsers).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbUsers)
}

func (r *userRepository) ListByRole(ctx context.Context, role string, limit, offset int) ([]*entity.User, error) {
	var dbUsers []model.User

	query := r.db.WithContext(ctx).Where("role = ?", role)
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&dbUsers).Error; err != nil {
		return nil, err
	}

	return r.modelsToEntities(dbUsers)
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, id int) error {
	now := time.Now()
	return r.db.WithContext(ctx).Model(&model.User{}).Where("id = ?", id).Update("last_login_at", now).Error
}

// Helper methods
func (r *userRepository) entityToModel(user *entity.User) *model.User {
	return &model.User{
		ID:            user.ID,
		Email:         user.Email,
		PasswordHash:  user.PasswordHash,
		Role:          user.Role.String(),
		IsActive:      user.IsActive,
		EmailVerified: user.EmailVerified,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		LastLoginAt:   user.LastLoginAt,
	}
}

func (r *userRepository) modelToEntity(dbUser *model.User) (*entity.User, error) {
	role, err := vo.ParseUserRole(dbUser.Role)
	if err != nil {
		return nil, err
	}

	return &entity.User{
		ID:            dbUser.ID,
		Email:         dbUser.Email,
		PasswordHash:  dbUser.PasswordHash,
		Role:          role,
		IsActive:      dbUser.IsActive,
		EmailVerified: dbUser.EmailVerified,
		CreatedAt:     dbUser.CreatedAt,
		UpdatedAt:     dbUser.UpdatedAt,
		LastLoginAt:   dbUser.LastLoginAt,
	}, nil
}

func (r *userRepository) modelsToEntities(dbUsers []model.User) ([]*entity.User, error) {
	entities := make([]*entity.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		entity, err := r.modelToEntity(&dbUser)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}
