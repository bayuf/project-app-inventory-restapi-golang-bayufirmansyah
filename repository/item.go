package repository

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"go.uber.org/zap"
)

type ItemRepository struct {
	DB     db.DBExecutor
	Logger *zap.Logger
}

func NewItemRepository(db db.DBExecutor, log *zap.Logger) *ItemRepository {
	return &ItemRepository{
		DB:     db,
		Logger: log,
	}
}
