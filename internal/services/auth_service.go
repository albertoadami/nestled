package services

import (
	"github.com/albertoadami/nestled/internal/auth"
	"github.com/albertoadami/nestled/internal/config"
	"github.com/albertoadami/nestled/internal/crypto"
	"github.com/albertoadami/nestled/internal/errors"
	"github.com/albertoadami/nestled/internal/repositories"
)

type AuthService interface {
	GenerateToken(username string, password string) (*auth.Token, error)
}

type authService struct {
	userRepository repositories.UserRepository
	JWtConfig      config.JWTConfig
}

func NewAuthService(userRepository repositories.UserRepository, jwtConfig config.JWTConfig) AuthService {
	return &authService{
		userRepository: userRepository,
		JWtConfig:      jwtConfig,
	}
}

func (s *authService) GenerateToken(username string, password string) (*auth.Token, error) {
	user, err := s.userRepository.GetUserByUsername(username)
	if err != nil {
		return nil, err
	}

	if user == nil {
		return nil, errors.CredentialsInvalid
	}

	valid := crypto.CheckPassword(password, user.PasswordHash)
	if valid {
		jwtToken, err := auth.GenerateToken(user.Id, s.JWtConfig.Secret, s.JWtConfig.Expiration)
		if err != nil {
			return nil, err
		}
		return jwtToken, nil

	}

	return nil, errors.CredentialsInvalid
}
