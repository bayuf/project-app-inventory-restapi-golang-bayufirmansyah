package repository

import (
	"context"
	"errors"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserRepository struct {
	DB     db.DBExecutor
	Logger *zap.Logger
}

func NewUserRepository(db db.DBExecutor, log *zap.Logger) *UserRepository {
	return &UserRepository{
		DB:     db,
		Logger: log,
	}
}

func (r *UserRepository) GetUserById(ctx context.Context, userId uuid.UUID) (*model.User, error) {
	query := `SELECT u.name, r.name
	FROM users u
	JOIN roles r ON u.role_id = r.id
	WHERE u.id=$1 AND u.deleted_at IS NULL AND u.is_active=true;`

	user := model.User{}
	if err := r.DB.QueryRow(ctx, query, userId).Scan(
		&user.ModelUser.Name,
		&user.ModelUser.RoleName,
	); err != nil {
		r.Logger.Error("cant scan get user id", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) AddUser(ctx context.Context, newUserData model.User) error {
	query := `INSERT INTO users (id, name, email, password_hash, role_id, is_active, deleted_at, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, TRUE, NULL, NOW(), NOW());`

	newUser := newUserData
	commandTag, err := r.DB.Exec(ctx, query,
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
