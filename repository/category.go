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

	return nil
}
