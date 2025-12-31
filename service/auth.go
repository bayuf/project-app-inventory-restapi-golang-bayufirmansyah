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

func (s *AuthService) Login(userReq dto.UserReq) (*dto.ResponseSession, error) {
	// get user data
	user, err := s.Repo.FindUserByEmail(model.User{Email: userReq.Email})
	if err != nil {
		s.Logger.Info("email not found", zap.Error(err))
		return nil, errors.New("invalid credentials")
	}

	// check password
	if !utils.CheckPassword(user.Password, userReq.Password) {
		s.Logger.Info("password not match", zap.Error(err))
		return nil, errors.New("invalid credentials")
	}

	// user active check
	if !user.IsActive {
		s.Logger.Info("user is inactive (banned or suspended)")
		return nil, errors.New("email is banned or suspended")
	}

	// revoke active session
	if err := s.Repo.RevokeSessionByUserId(user.ModelUser.ID); err != nil {
		return nil, err
	}

	// create new session
	if err := CreateSession(s, user.ModelUser.ID); err != nil {
		return nil, err
	}

	// get session ID
	session, err := GetSession(s, user.ModelUser.ID)
	return session, nil
}

func (s *AuthService) Logout(sessionId uuid.UUID) error {
	if err := s.Repo.RevokeSessionById(sessionId); err != nil {
		return err
	}
	return nil
}

func (s *AuthService) ValidateSession(sessionId uuid.UUID) (*dto.ResponseSession, error) {
	return s.Repo.ValidateSession(sessionId)
}

func CreateSession(s *AuthService, userId uuid.UUID) error {
	err := s.Repo.CreateSession(dto.Session{
		ID:        uuid.New(),
		UserID:    userId,
		ExpiresAt: time.Now().Add(24 * time.Hour),
	})
	if err != nil {
		return err
	}
	return nil
}

func GetSession(s *AuthService, userId uuid.UUID) (*dto.ResponseSession, error) {
	session, err := s.Repo.GetSessionByUserId(userId)
	if err != nil {
		return nil, err
	}
	return session, nil
}
