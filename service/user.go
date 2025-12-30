package service

import (
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

func (s *UserService) AddUser(newUserData dto.UserAdd) error {
	hashedPassword, err := utils.HashPassword(newUserData.Password)
	if err != nil {
		s.Logger.Error("hashing failed")
		return err
	}
	// call repo and send model user
	if err := s.Repo.AddUser(model.User{
		ID:       uuid.New(), // create uuid
		Name:     newUserData.Name,
		Email:    newUserData.Email,
		Password: hashedPassword,
		RoleID:   newUserData.RoleID,
	}); err != nil {
		return err
	}
	return nil
}
