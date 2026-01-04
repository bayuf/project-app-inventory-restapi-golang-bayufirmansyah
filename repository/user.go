package repository

import (
	"context"
	"errors"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/db"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
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
	query := `SELECT u.name, r.name, u.email, u.created_at
	FROM users u
	JOIN roles r ON u.role_id = r.id
	WHERE u.id=$1 AND u.deleted_at IS NULL AND u.is_active=true;`

	user := model.User{}
	if err := r.DB.QueryRow(ctx, query, userId).Scan(
		&user.ModelUser.Name,
		&user.ModelUser.RoleName,
		&user.Email,
		&user.ModelUser.Created_At,
	); err != nil {
		r.Logger.Error("cant scan get user id", zap.Error(err))
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetAllUsers(ctx context.Context, page, limit int) (*[]model.User, int, error) {
	offset := (page - 1) * limit

	// get total data for pagination
	var total int
	countQuery := `SELECT COUNT(*) FROM users WHERE deleted_at IS NULL AND is_active = TRUE`
	err := r.DB.QueryRow(ctx, countQuery).Scan(&total)
	if err != nil {
		r.Logger.Error("error query getAll repo ", zap.Error(err))
		return nil, 0, err
	}

	query := `SELECT u.name, r.name, u.email, u.created_at
	FROM users u
	JOIN roles r ON u.role_id = r.id
	WHERE u.deleted_at IS NULL AND u.is_active=true
	ORDER BY u.created_at
	LIMIT $1 OFFSET $2;`

	rows, err := r.DB.Query(ctx, query, limit, offset)
	if err != nil {
		r.Logger.Error("failed get all users repository", zap.Error(err))
		return nil, 0, err
	}
	defer rows.Close()

	users := []model.User{}
	for rows.Next() {
		user := model.User{}
		rows.Scan(&user.ModelUser.Name, &user.ModelUser.RoleName, &user.Email, &user.ModelUser.Created_At)

		users = append(users, user)
	}

	return &users, total, nil
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

func (r *UserRepository) RegisterUser(ctx context.Context, newUserData model.User) error {
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

	r.Logger.Info("new user registered", zap.Any("data", newUserData))
	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	query := `
	UPDATE users
	SET updated_at = NOW(), deleted_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL;`

	commandTag, err := r.DB.Exec(ctx, query, userId)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Error("user not found", zap.Error(err))
			return err
		}
		r.Logger.Error("cant delete user", zap.Error(err))
		return err
	}

	r.Logger.Info("user deleted", zap.Any("ID", userId))
	return nil
}

func (r *UserRepository) SuspendUser(ctx context.Context, userData dto.UserSuspend) error {
	query := `
	UPDATE users
	SET is_active = $2, updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL AND is_active = TRUE;`

	commandTag, err := r.DB.Exec(ctx, query, userData.ID, userData.Suspend)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Error("user not found", zap.Error(err))
			return err
		}
		r.Logger.Error("cant suspend user", zap.Error(err))
		return err
	}

	r.Logger.Info("user suspended", zap.Any("ID", userData.ID))
	return nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, newData dto.UserUpdate) error {
	query := `
	UPDATE users
	SET name = $2, password_hash = $3, role_id = $4, updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL AND is_active = TRUE;`

	commandTag, err := r.DB.Exec(ctx, query, newData.ID, newData.Name, newData.Password, newData.RoleID)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Error("user not found", zap.Error(err))
			return err
		}
		r.Logger.Error("cant update user data", zap.Error(err))
		return err
	}

	r.Logger.Info("user data updated", zap.Any("ID", newData.ID))
	return nil
}

func (r *UserRepository) UpdateMyUserData(ctx context.Context, newData dto.UserSelfUpdate) error {
	query := `
	UPDATE users
	SET name = $2, password_hash = $3, updated_at = NOW()
	WHERE id = $1 AND deleted_at IS NULL AND is_active = TRUE;`

	commandTag, err := r.DB.Exec(ctx, query, newData.ID, newData.Name, newData.Password)
	if err != nil {
		if commandTag.RowsAffected() == 0 {
			r.Logger.Error("user not found", zap.Error(err))
			return err
		}
		r.Logger.Error("cant update user data", zap.Error(err))
		return err
	}

	r.Logger.Info("user data updated", zap.Any("ID", newData.ID))
	return nil
}
