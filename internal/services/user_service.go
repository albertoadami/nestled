package services

import (
	"github.com/albertoadami/nestled/internal/crypto"
	"github.com/albertoadami/nestled/internal/dto"
	"github.com/albertoadami/nestled/internal/model"
	"github.com/albertoadami/nestled/internal/repositories"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(request *dto.CreateUserRequest) (uuid.UUID, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) UserService {
	return &userService{
		userRepository: userRepository,
	}
}

func (s *userService) CreateUser(request *dto.CreateUserRequest) (uuid.UUID, error) {

	hashedPassword, err := crypto.HashPassword(request.Password)
	if err != nil {
		return uuid.Nil, err
	}

	user := &model.User{
		Id:           uuid.New(),
		FirstName:    request.FirstName,
		LastName:     request.LastName,
		Username:     request.Username,
		Email:        request.Email,
		PasswordHash: hashedPassword,
		Status:       model.UserStatusPending,
	}

	return s.userRepository.CreateUser(user)

}
