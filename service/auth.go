package service

import (
	"errors"
	"time"

	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/dto"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/model"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/repository"
	"github.com/bayuf/project-app-inventory-restapi-golang-bayufirmansyah/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type AuthService struct {
	Repo   *repository.AuthRepository
	Logger *zap.Logger
}

func NewAuthService(repo *repository.AuthRepository, log *zap.Logger) *AuthService {
	return &AuthService{
		Repo:   repo,
		Logger: log,
	}
}

func (s *AuthService) Login(userReq dto.UserReq) (dto.Session, error) {
	// get user data
	user, err := s.Repo.FindUserByEmail(model.User{Email: userReq.Email})
	if err != nil {
		s.Logger.Info("email not found", zap.Error(err))
		return dto.Session{}, errors.New("invalid credentials")
	}

	// check password
	if !utils.CheckPassword(user.Password, userReq.Password) {
		s.Logger.Info("password not match", zap.Error(err))
		return dto.Session{}, errors.New("invalid credentials")
	}

	// user active check
	if !user.IsActive {
		s.Logger.Info("user is inactive (banned)")
		return dto.Session{}, errors.New("email is banned")
	}

	// create sessions ID
	session, err := s.Repo.CreateSession(dto.Session{
		ID:        uuid.New(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(15 * time.Minute),
	})
	if err != nil {
		return dto.Session{}, err
	}

	return session, nil
}
