package handlers

import (
	"fmt"
	"net/http"

	"github.com/albertoadami/nestled/internal/dto"
	"github.com/albertoadami/nestled/internal/errors"
	"github.com/albertoadami/nestled/internal/services"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UserHandler struct {
	userService services.UserService
	logger      *zap.Logger
}

func NewUserHandler(userService services.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

func (u *UserHandler) RegisterUser(c *gin.Context) {

	var request dto.CreateUserRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	userId, err := u.userService.CreateUser(&request)
	if err != nil {
		switch err {
		case errors.ErrUsernameAlreadyExists:
			c.JSON(http.StatusConflict, dto.NewErrorResponse(err.Error(), fmt.Sprintf("The username %s is already in use", request.Username)))
			return
		case errors.ErrEmailAlreadyExists:
			c.JSON(http.StatusConflict, dto.NewErrorResponse(err.Error(), fmt.Sprintf("The email %s is already in use", request.Email)))
			return
		case errors.ErrPasswordTooWeak:
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		default:
			c.Status(500)
			return
		}
	}

	u.logger.Info(fmt.Sprintf("User created with ID: %s", userId))

	locationPathResponse := fmt.Sprintf("/api/v1/users/%s", userId)
	c.Header("Location", locationPathResponse)
	c.Status(http.StatusCreated)

}

func (u *UserHandler) GetCurrentUser(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Status(http.StatusInternalServerError)
		return
	}
	userIdUUID, err := uuid.Parse(userId.(string))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	user, err := u.userService.GetUserById(userIdUUID)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	if user == nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	response := &dto.UserResponse{
		Id:        user.Id.String(),
		Username:  user.Username,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	}
	c.JSON(http.StatusOK, response)

}
