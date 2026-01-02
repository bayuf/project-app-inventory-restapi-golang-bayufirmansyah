package repository

import (
	"context"
	"errors"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"go.uber.org/zap"
)

type CategoryRepository struct {
	DB     db.DBExecutor
	Logger *zap.Logger
}

func NewCategoryRepository(db db.DBExecutor, log *zap.Logger) *CategoryRepository {
	return &CategoryRepository{
		DB:     db,
		Logger: log,
	}
}

func (r *CategoryRepository) Create(ctx context.Context, newData model.Category) error {
	query := `INSERT INTO categories (name, description)
	VALUES ($1, $2);`

	commandTag, err := r.DB.Exec(ctx, query, newData.Name, newData.Description)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Error("no row affected to create new category", zap.Error(err))
			return errors.New("no row affected")
		}
		r.Logger.Error("cant create new category ", zap.Error(err))
		return err
	}

	r.Logger.Info("new category created with ", zap.Int("ID", newData.ID))
	return nil
}

func (r *CategoryRepository) GetAll(ctx context.Context, page, limit int) (*[]model.Category, int, error) {
	offset := (page - 1) * limit

	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM categories WHERE deleted_at IS NULL`
	err := r.DB.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error query getAll repo ", zap.Error(err))
		return nil, 0, err
	}

	query := `SELECT id, name, description, created_at, updated_at
	FROM categories
	WHERE deleted_at IS NULL
	ORDER BY id ASC
	LIMIT $1 OFFSET $2;`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		if rows.CommandTag().RowsAffected() == 0 {
			return nil, 0, errors.New("categories is empty")
		}

		r.Logger.Error("error get all categories ", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	categories := []model.Category{}
	for rows.Next() {
		category := model.Category{}
		rows.Scan(
			&category.ID, &category.Name, &category.Description, &category.Created_At, &category.Updated_At,
		)
		categories = append(categories, category)
	}

	return &categories, total, nil
}

func (r *CategoryRepository) GetById(ctx context.Context, id int) (*model.Category, error) {
	query := `SELECT id, name, description, created_at, updated_at
	FROM categories
	WHERE id = $1 AND deleted_at IS NULL;`

	category := model.Category{}
	if err := r.DB.QueryRow(ctx, query, id).Scan(
		&category.ID, &category.Name, &category.Description, &category.Created_At, &category.Updated_At,
	); err != nil {
		r.Logger.Error("error getById category", zap.Error(err))
		return nil, err
	}

	return &category, nil
}

func (r *CategoryRepository) Update(ctx context.Context, new model.Category) error {
	query := `UPDATE categories SET name = $2, description = $3, updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL;`

	commandTag, err := r.DB.Exec(ctx, query, new.ID, new.Name, new.Description)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Error("category cant be updated", zap.Error(err))
			return err
		}
		r.Logger.Error("cant update category", zap.Error(err))
		return err
	}

	r.Logger.Info("updated category", zap.Int("ID", new.ID))
	return nil
}

func (r *CategoryRepository) Delete(ctx context.Context, id int) error {
	query := `UPDATE categories SET deleted_at = NOW(), updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL;`

	commandTag, err := r.DB.Exec(ctx, query, id)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Error("category cant be deleted", zap.Error(err))
			return err
		}
		r.Logger.Error("cant delete category", zap.Error(err))
		return err
	}

	r.Logger.Info("deleted category", zap.Int("ID", id))
	return nil
}
