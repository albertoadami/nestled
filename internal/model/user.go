package model

import (
	"github.com/google/uuid"
)

type UserStatus string

const (
	UserStatusActive  UserStatus = "ACTIVE"
	UserStatusBlocked UserStatus = "BLOCKED"
	UserStatusPending UserStatus = "PENDING"
)

type User struct {
	Id           uuid.UUID
	Username     string
	FirstName    string
	LastName     string
	Email        string
	PasswordHash string
	Status       UserStatus
}
