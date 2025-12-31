package repository

import (
	"context"
	"errors"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type AuthRepository struct {
	DB     db.DBExecutor
	Logger *zap.Logger
}

func NewAuthRepository(db db.DBExecutor, log *zap.Logger) *AuthRepository {
	return &AuthRepository{
		DB:     db,
		Logger: log,
	}
}

func (r *AuthRepository) FindUserByEmail(ctx context.Context, userReq model.User) (*model.User, error) {
	query := `SELECT id, name, email, password_hash, role_id, is_active, created_at, updated_at
	FROM users
	WHERE email=$1 AND deleted_at IS NULL;`

	row := r.DB.QueryRow(ctx, query, userReq.Email)

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
			return &model.User{}, errors.New("user not found")
		}

		r.Logger.Error("failed scan user data", zap.Error(err))
		return &model.User{}, err
	}

	return &user, nil
}

func (r *AuthRepository) GetSessionByUserId(ctx context.Context, userId uuid.UUID) (*dto.ResponseSession, error) {
	query := `SELECT s.id, s.user_id, u.name, u.role_id, r.name, s.created_at, s.expires_at
	FROM sessions s
	JOIN users u ON u.id = s.user_id
	JOIN roles r ON u.role_id = r.id
	WHERE s.user_id=$1 AND s.revoked_at IS NULL AND s.expires_at > NOW();`
	session := dto.ResponseSession{}
	if err := r.DB.QueryRow(ctx, query, userId).
		Scan(
			&session.ID,
			&session.UserID,
			&session.Username,
			&session.RoleId,
			&session.RoleName,
			&session.CreatedAt,
			&session.ExpiresAt,
		); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.Logger.Info("session not found")
			return nil, errors.New("session not found")
		}
		r.Logger.Error("cant scan row check session status", zap.Error(err))
		return nil, err
	}

	return &session, nil
}

func (r *AuthRepository) CreateSession(ctx context.Context, req dto.Session) error {
	query := `INSERT INTO sessions (id, user_id, expires_at)
	VALUES ($1, $2, $3);`

	commandTag, err := r.DB.Exec(ctx, query, req.ID, req.UserID, req.ExpiresAt)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Error("failed create new sessions", zap.String("error:", "0 rows affected"))
			return err
		}

		r.Logger.Error("failed create new sessions", zap.Error(err))
		return err
	}

	return nil
}

func (r *AuthRepository) RevokeSessionByUserId(ctx context.Context, userId uuid.UUID) error {
	query := `UPDATE sessions SET revoked_at = NOW()
			WHERE user_id = $1 AND revoked_at IS NULL AND expires_at > NOW();`

	_, err := r.DB.Exec(ctx, query, userId)
	if err != nil {
		r.Logger.Error("cant revoked sessions User ID")
		return err
	}

	return nil
}

func (r *AuthRepository) RevokeSessionById(ctx context.Context, sessionId uuid.UUID) error {
	query := `UPDATE sessions SET revoked_at = NOW()
			WHERE id = $1 AND revoked_at IS NULL AND expires_at > NOW();`

	_, err := r.DB.Exec(ctx, query, sessionId)
	if err != nil {
		r.Logger.Error("cant revoked sessions by ID")
		return err
	}

	return nil
}

func (r *AuthRepository) ValidateSession(ctx context.Context, sessionId uuid.UUID) (*dto.ResponseSession, error) {
	query := `SELECT s.id, s.user_id, u.name, r.name
			FROM sessions s
			JOIN users u ON u.id = s.user_id
			JOIN roles r ON u.role_id = r.id
			WHERE s.id = $1
			  AND s.revoked_at IS NULL
			  AND s.expires_at > NOW()
			  AND u.is_active = TRUE
			  AND u.deleted_at IS NULL;`
	session := dto.ResponseSession{}
	if err := r.DB.QueryRow(context.Background(), query, sessionId).Scan(
		&session.ID,
		&session.UserID,
		&session.Username,
		&session.RoleName,
	); err != nil {
		r.Logger.Error("error validate session", zap.Error(err))
		return nil, err
	}
	return &session, nil
}
