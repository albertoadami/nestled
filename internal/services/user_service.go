package services

import (
	"github.com/albertoadami/nestled/internal/crypto"
	"github.com/albertoadami/nestled/internal/dto"
	"github.com/albertoadami/nestled/internal/errors"
	"github.com/albertoadami/nestled/internal/model"
	"github.com/albertoadami/nestled/internal/repositories"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(request *dto.CreateUserRequest) (uuid.UUID, error)
	GetUserById(id uuid.UUID) (*model.User, error)
	ChangePassword(user *model.User, currentPassword string, newPassword string) error
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

	return s.userRepository.Create(user)

}

func (s *userService) GetUserById(id uuid.UUID) (*model.User, error) {
	return s.userRepository.GetUserById(id)
}

func (s *userService) ChangePassword(user *model.User, currentPassword string, newPassword string) error {
	if !crypto.CheckPassword(currentPassword, user.PasswordHash) {
		return errors.ErrInvalidPassword
	}

	// generate the new hash for the new password
	hashedPassword, err := crypto.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword

	return s.userRepository.Update(user)

}
