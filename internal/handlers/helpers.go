package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func getUserIdFromContext(c *gin.Context) (uuid.UUID, bool) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Status(http.StatusInternalServerError)
		return uuid.Nil, false
	}
	userIdUUID, err := uuid.Parse(userId.(string))
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return uuid.Nil, false
	}
	return userIdUUID, true
}
