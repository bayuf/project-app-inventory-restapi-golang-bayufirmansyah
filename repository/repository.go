package repository

import (
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"go.uber.org/zap"
)

type Repository struct {
}

func NewRepository(db db.PgxIface, log *zap.Logger) *Repository {
	return &Repository{}
}
