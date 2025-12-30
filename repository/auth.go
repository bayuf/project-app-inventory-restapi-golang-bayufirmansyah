package repository

import (
	"context"
	"errors"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type AuthRepository struct {
	DB     db.PgxIface
	Logger *zap.Logger
}

func NewAuthRepository(db db.PgxIface, log *zap.Logger) *AuthRepository {
	return &AuthRepository{
		DB:     db,
		Logger: log,
	}
}

func (r *AuthRepository) FindUserByEmail(userReq model.User) (model.User, error) {
	query := `SELECT id, name, email, password_hash, role_id, is_active, created_at, updated_at
	FROM users
	WHERE email=$1 AND deleted_at IS NULL;`

	row := r.DB.QueryRow(context.Background(), query, userReq.Email)

	user := model.User{}
	if err := row.Scan(
		&user.ModelUser.ID,
		&user.ModelUser.Name,
		&user.Email,
		&user.Password,
		&user.RoleID,
		&user.ModelUser.IsActive,
		&user.ModelUser.Created_At,
		&user.ModelUser.Updated_At,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.Logger.Info("user not found", zap.String("email :", userReq.Email))
			return model.User{}, errors.New("user not found")
		}

		r.Logger.Error("failed scan user data", zap.Error(err))
		return model.User{}, err
	}

	return user, nil
}

func (r *AuthRepository) CreateSession(req dto.Session) (dto.Session, error) {
	query := `INSERT INTO sessions (id, user_id, expires_at)
	VALUES ($1, $2, $3);`

	commandTag, err := r.DB.Exec(context.Background(), query, req.ID, req.UserID, req.ExpiresAt)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Error("failed create new sessions", zap.String("error:", "0 rows affected"))
			return dto.Session{}, err
		}

		r.Logger.Error("failed create new sessions", zap.Error(err))
		return dto.Session{}, err
	}

	return req, nil
}
