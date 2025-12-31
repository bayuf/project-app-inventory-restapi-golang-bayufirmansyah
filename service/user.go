package service

import (
	"context"

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
		Name:     user.ModelUser.Name,
		RoleName: user.ModelUser.RoleName,
	}, nil
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
