package middleware

import (
	"net/http"
	"strings"

	"github.com/albertoadami/nestled/internal/auth"
	"github.com/albertoadami/nestled/internal/dto"
	"github.com/gin-gonic/gin"
)

func BearerAuthentication(tokenManager *auth.TokenManager) gin.HandlerFunc {

	return func(c *gin.Context) {

		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.ErrorResponse{Message: "missing or invalid Authorization header"})
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		tokenInfo, err := tokenManager.ParseToken(token)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, &dto.ErrorResponse{Message: "invalid token"})
			return
		}

		c.Set("userId", tokenInfo.UserId)
		c.Next()

	}

}
