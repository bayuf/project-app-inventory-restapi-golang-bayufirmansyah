package service

import (
	"context"
	"errors"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/repository"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserService struct {
	Repo   *repository.UserRepository
	Logger *zap.Logger
}

func NewUserService(repo *repository.UserRepository, log *zap.Logger) *UserService {
	return &UserService{
		Repo:   repo,
		Logger: log,
	}
}

func (s *UserService) GetUserData(ctx context.Context, userId uuid.UUID) (*dto.UserResponse, error) {
	user, err := s.Repo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return &dto.UserResponse{
		Name:       user.ModelUser.Name,
		RoleName:   user.ModelUser.RoleName,
		Email:      user.Email,
		Created_At: user.Created_At,
	}, nil
}

func (s *UserService) GetAllUsersData(ctx context.Context, page, limit int) (*[]dto.UserResponse, *dto.Pagination, error) {
	users, total, err := s.Repo.GetAllUsers(ctx, page, limit)
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.Pagination{
		CurrentPage:  page,
		Limit:        limit,
		TotalPages:   utils.TotalPage(limit, int64(total)),
		TotalRecords: total,
	}

	usersRes := []dto.UserResponse{}
	for _, v := range *users {
		res := dto.UserResponse{
			Name:       v.ModelUser.Name,
			RoleName:   v.ModelUser.RoleName,
			Email:      v.Email,
			Created_At: v.ModelUser.Created_At,
		}

		usersRes = append(usersRes, res)
	}

	return &usersRes, &pagination, nil
}

func (s *UserService) AddUser(ctx context.Context, newUserData dto.UserAdd) error {
	hashedPassword, err := utils.HashPassword(newUserData.Password)
	if err != nil {
		s.Logger.Error("hashing failed")
		return err
	}
	// call repo and send model user
	if err := s.Repo.AddUser(ctx, model.User{
		ModelUser: model.ModelUser{
			ID:   uuid.New(),
			Name: newUserData.Name,
		},
		Email:    newUserData.Email,
		Password: hashedPassword,
		RoleID:   newUserData.RoleID,
	}); err != nil {
		return err
	}
	return nil
}

func (s *UserService) RegisterUser(ctx context.Context, newUserData dto.UserRegister) error {
	// hash string password
	hashedPassword, err := utils.HashPassword(newUserData.Password)
	if err != nil {
		s.Logger.Error("hashing failed")
		return err
	}
	// call repo and send model user
	if err := s.Repo.RegisterUser(ctx, model.User{
		ModelUser: model.ModelUser{
			ID:   uuid.New(),
			Name: newUserData.Name,
		},
		Email:    newUserData.Email,
		Password: hashedPassword,
		RoleID:   3, //user register as staff
	}); err != nil {
		return err
	}
	return nil
}

func (s *UserService) DeleteUser(ctx context.Context, userId uuid.UUID) error {
	return s.Repo.DeleteUser(ctx, userId)
}

func (s *UserService) SuspendUser(ctx context.Context, userData dto.UserSuspend) error {
	return s.Repo.SuspendUser(ctx, dto.UserSuspend{
		ID:      userData.ID,
		Suspend: !userData.Suspend, // status user is not active (convert true to false and vice versa)
	})
}

func (s *UserService) UpdateUser(ctx context.Context, userData dto.UserUpdate) error {
	newHashedPass, err := utils.HashPassword(userData.Password)
	if err != nil {
		return err
	}

	// if try change to super admin
	superAdmin := 1
	if userData.RoleID == superAdmin {
		return errors.New("cant change to super admin")
	}

	return s.Repo.UpdateUser(ctx, dto.UserUpdate{
		ID:       userData.ID,
		Name:     userData.Name,
		Password: newHashedPass,
		RoleID:   userData.RoleID,
	})
}

func (s *UserService) UpdateMyUserData(ctx context.Context, userData dto.UserSelfUpdate) error {
	newHashedPass, err := utils.HashPassword(userData.Password)
	if err != nil {
		return err
	}

	return s.Repo.UpdateMyUserData(ctx, dto.UserSelfUpdate{
		ID:       userData.ID,
		Name:     userData.Name,
		Password: newHashedPass,
	})
}
