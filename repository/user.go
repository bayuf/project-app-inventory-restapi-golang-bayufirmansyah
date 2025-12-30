package repository

import (
	"context"
	"errors"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"go.uber.org/zap"
)

type UserRepository struct {
	DB     db.PgxIface
	Logger *zap.Logger
}

func NewUserRepository(db db.PgxIface, log *zap.Logger) *UserRepository {
	return &UserRepository{
		DB:     db,
		Logger: log,
	}
}

func (r *UserRepository) AddUser(newUserData model.User) error {
	query := `INSERT INTO users (id, name, email, password_hash, role_id, is_active, deleted_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, TRUE, NULL, NOW(), NOW());`

	newUser := newUserData
	commandTag, err := r.DB.Exec(context.Background(), query,
		newUser.ModelUser.ID,
		newUser.ModelUser.Name,
		newUser.Email,
		newUser.Password,
		newUser.RoleID,
	)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Error("0 row affected in AddUser repository", zap.Error(err))
			return errors.New("failed insert data")
		}

		r.Logger.Error("error insert new user repository", zap.Error(err))
		return err
	}

	r.Logger.Info("add new user", zap.Any("data", newUserData))
	return nil
}
