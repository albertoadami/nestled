package testhelpers

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func MockAuthentication(userId uuid.UUID) gin.HandlerFunc {

	return func(c *gin.Context) {

		c.Set("userId", userId.String())
		c.Next()

	}
}
