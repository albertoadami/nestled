package repositories

import (
	"context"
	"database/sql"
	"errors"

	customErrors "github.com/albertoadami/nestled/internal/errors"
	"github.com/albertoadami/nestled/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type ActivationTokenRepository interface {
	Create(activationToken *model.ActivationToken) (uuid.UUID, error)
	GetByToken(token string) (*model.ActivationToken, error)
	DeleteById(id uuid.UUID) error
}

type activationTokenRepository struct {
	db *sqlx.DB
}

func NewActivationTokenRepository(db *sqlx.DB) ActivationTokenRepository {
	return &activationTokenRepository{db: db}
}

func (r *activationTokenRepository) Create(activationToken *model.ActivationToken) (uuid.UUID, error) {
	query := `INSERT INTO activation_tokens (id, user_id, token, expires_at)
              VALUES ($1, $2, $3, $4)
              RETURNING id`
	var id uuid.UUID
	err := r.db.QueryRowx(query,
		activationToken.Id,
		activationToken.UserId,
		activationToken.Token,
		activationToken.ExpiresAt,
	).Scan(&id)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil

}

func (r *activationTokenRepository) GetByToken(token string) (*model.ActivationToken, error) {
	query := `SELECT id, user_id, token, expires_at
			 FROM activation_tokens
			 WHERE token = $1`

	var result model.ActivationToken
	err := r.db.Get(&result, query, token)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &result, nil
}

func (r *activationTokenRepository) DeleteById(id uuid.UUID) error {
	query := `DELETE FROM activation_tokens WHERE id = $1`

	// TODO: refactor the code to work with Context
	result, err := r.db.ExecContext(context.Background(), query, id)
	if err != nil {
		return err
	}

	affectedRows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affectedRows == 0 {
		return customErrors.ErrNoRowsAffected
	}
	return nil

}
