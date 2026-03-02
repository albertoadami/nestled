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
	Id           uuid.UUID  `db:"id"`
	Username     string     `db:"username"`
	FirstName    string     `db:"first_name"`
	LastName     string     `db:"last_name"`
	Email        string     `db:"email"`
	PasswordHash string     `db:"password_hash"`
	Status       UserStatus `db:"status"`
}
