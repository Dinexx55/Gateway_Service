package service

import (
	"errors"
	"go.uber.org/zap"
)

type UserRepository interface {
	GetUserByLogin(login string) (*User, error)
}

type AuthProvider interface {
	GetJWTToken(login string) (string, error)
}

type User struct {
	Login    string
	Password string
}

type AuthService struct {
	provider   AuthProvider
	logger     *zap.Logger
	repository UserRepository
}

func NewAuthService(provider AuthProvider, logger *zap.Logger, repository UserRepository) *AuthService {
	return &AuthService{
		provider:   provider,
		logger:     logger,
		repository: repository,
	}
}

var (
	ErrUserNotFound    = errors.New("user with provided login does not exist")
	ErrInvalidPassword = errors.New("invalid password for user")
)

func (s *AuthService) SignIn(credentials User) (string, error) {
	user, err := s.repository.GetUserByLogin(credentials.Login)
	if user == nil {
		return "", ErrUserNotFound
	}

	isPasswordValid, err := checkPassword(credentials.Password, user.Password)
	if !isPasswordValid {
		return "", ErrInvalidPassword
	}

	accessToken, err := s.provider.GetJWTToken(credentials.Login)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

// there should be some hash logic implementation
func checkPassword(password, hash string) (bool, error) {
	if password == hash {
		return true, nil
	} else {
		return false, nil
	}
}
