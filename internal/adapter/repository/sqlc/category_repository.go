package repository

import (
	"context"
	"database/sql"
	"errors"

	sqlc "github.com/hydr0g3nz/poc_pos_restuarant/internal/adapter/sqlc/generated"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/entity"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/repository"
	"github.com/hydr0g3nz/poc_pos_restuarant/internal/domain/vo"
	"github.com/jackc/pgx/v5/pgxpool"
)

type categoryRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewCategoryRepository(db *pgxpool.Pool) repository.CategoryRepository {
	return &categoryRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *categoryRepository) Create(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	dbCategory, err := r.queries.CreateCategory(ctx, category.Name.String())
	if err != nil {
		return nil, err
	}

	return r.dbCategoryToEntity(dbCategory)
}

func (r *categoryRepository) GetByID(ctx context.Context, id int) (*entity.Category, error) {
	dbCategory, err := r.queries.GetCategoryByID(ctx, int32(id))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.dbCategoryToEntity(dbCategory)
}

func (r *categoryRepository) GetByName(ctx context.Context, name string) (*entity.Category, error) {
	dbCategory, err := r.queries.GetCategoryByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return r.dbCategoryToEntity(dbCategory)
}

func (r *categoryRepository) Update(ctx context.Context, category *entity.Category) (*entity.Category, error) {
	dbCategory, err := r.queries.UpdateCategory(ctx, sqlc.UpdateCategoryParams{
		ID:   int32(category.ID),
		Name: category.Name.String(),
	})
	if err != nil {
		return nil, err
	}

	return r.dbCategoryToEntity(dbCategory)
}

func (r *categoryRepository) Delete(ctx context.Context, id int) error {
	return r.queries.DeleteCategory(ctx, int32(id))
}

func (r *categoryRepository) List(ctx context.Context) ([]*entity.Category, error) {
	dbCategories, err := r.queries.ListCategories(ctx)
	if err != nil {
		return nil, err
	}

	return r.dbCategoriesToEntities(dbCategories)
}

func (r *categoryRepository) HasMenuItems(ctx context.Context, categoryID int) (bool, error) {
	result, err := r.queries.CheckCategoryHasMenuItems(ctx, int32(categoryID))
	if err != nil {
		return false, err
	}
	return result, nil
}

// Helper methods for conversion
func (r *categoryRepository) dbCategoryToEntity(dbCategory *sqlc.Category) (*entity.Category, error) {
	categoryType, err := vo.NewCategoryType(dbCategory.Name)
	if err != nil {
		return nil, err
	}

	return &entity.Category{
		ID:        int(dbCategory.ID),
		Name:      categoryType,
		CreatedAt: dbCategory.CreatedAt.Time,
	}, nil
}

func (r *categoryRepository) dbCategoriesToEntities(dbCategories []*sqlc.Category) ([]*entity.Category, error) {
	entities := make([]*entity.Category, len(dbCategories))
	for i, dbCategory := range dbCategories {
		entity, err := r.dbCategoryToEntity(dbCategory)
		if err != nil {
			return nil, err
		}
		entities[i] = entity
	}
	return entities, nil
}
