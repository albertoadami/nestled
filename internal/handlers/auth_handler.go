package handlers

import (
	"errors"

	"github.com/albertoadami/nestled/internal/dto"
	appErrorrs "github.com/albertoadami/nestled/internal/errors"
	"github.com/albertoadami/nestled/internal/services"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

func (h *AuthHandler) GenerateToken(c *gin.Context) {
	var request dto.LoginRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	token, err := h.authService.GenerateToken(request.Username, request.Password)

	switch {
	case errors.Is(err, appErrorrs.CredentialsInvalid):
		c.JSON(401, &dto.ErrorResponse{
			Message: err.Error(),
			Details: "Invalid username or password",
		})
	default:
		tokenResponse := &dto.TokenResponse{
			Token:          token.Value,
			ExpirationTime: token.ExpirationTime.UnixMilli(),
		}
		c.JSON(200, tokenResponse)
	}

}
