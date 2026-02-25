package repositories

import (
	"github.com/albertoadami/nestled/internal/errors"
	"github.com/albertoadami/nestled/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type UserRepository interface {
	CreateUser(user *model.User) (uuid.UUID, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(user *model.User) (uuid.UUID, error) {
	query := `INSERT INTO users (id, username, first_name, last_name, email, password_hash, status)
              VALUES ($1, $2, $3, $4, $5, $6, $7)
              RETURNING id`

	var id uuid.UUID
	err := r.db.QueryRowx(query,
		user.Id,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.Status,
	).Scan(&id)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" {
			switch pqErr.Constraint {
			case "users_username_key":
				return uuid.Nil, errors.ErrUsernameAlreadyExists
			case "users_email_key":
				return uuid.Nil, errors.ErrEmailAlreadyExists
			}
		}
		return uuid.Nil, err
	}

	return id, nil
}
