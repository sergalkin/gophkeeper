package postgres

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v4"

	"github.com/sergalkin/gophkeeper/internal/server/model"
	"github.com/sergalkin/gophkeeper/internal/server/storage"
)

var _ storage.SecretServerStorage = (*SecretPostgresStorage)(nil)

type SecretPostgresStorage struct {
	conn *pgx.Conn
}

const (
	CreateSecrete = `
					insert into secrets (user_id, type_id, title, content, created_at, updated_at) 
					values ($1,$2,$3,$4,$5,$6) 
					returning id
`
	GetSecret = `select id, user_id, type_id, title, content, created_at, updated_at, deleted_at 
				 from secrets 
				 where id = $1 and user_id = $2
`
	DeleteSecret = `delete from secrets where id = $1 and user_id = $2`
)

func NewSecretPostgresStorage(c *pgx.Conn) *SecretPostgresStorage {
	return &SecretPostgresStorage{conn: c}
}

// CreateSecret - stores provided model.Secret in database.
//
// Values of key Content of model.Secret is being hex encoded.
func (s *SecretPostgresStorage) CreateSecret(ctx context.Context, secret model.Secret) (model.Secret, error) {
	ctxWithTimeOut, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := s.conn.QueryRow(ctxWithTimeOut, CreateSecrete, secret.UserID, secret.TypeID, secret.Title,
		hex.EncodeToString(secret.Content), secret.CreatedAt, secret.UpdatedAt,
	).Scan(&secret.ID)
	if err != nil {
		return secret, fmt.Errorf("error in storing secret in db: %w", err)
	}

	return secret, nil
}

// GetSecret - return rehydrated model.Secret from database.
//
// Searches by user_id and id from provided model.Secret
func (s *SecretPostgresStorage) GetSecret(ctx context.Context, secret model.Secret) (model.Secret, error) {
	ctxWithTimeOut, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := s.conn.QueryRow(ctxWithTimeOut, GetSecret, secret.ID, secret.UserID).Scan(
		&secret.ID, &secret.UserID, &secret.TypeID, &secret.Title, &secret.Content,
		&secret.CreatedAt, &secret.UpdatedAt, &secret.DeletedAt,
	)
	if err != nil {
		return secret, fmt.Errorf("error in getting secret from db: %w", err)
	}

	decode, decErr := hex.DecodeString(string(secret.Content))
	if decErr != nil {
		return secret, fmt.Errorf("error in decoding content from db: %w", decErr)
	}

	secret.Content = decode

	return secret, nil
}

// DeleteSecret - deletes a user secret from database.
func (s *SecretPostgresStorage) DeleteSecret(ctx context.Context, secret model.Secret) (model.Secret, error) {
	ctxWithTimeOut, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := s.conn.QueryRow(ctxWithTimeOut, DeleteSecret, secret.ID, secret.UserID).Scan()
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return secret, fmt.Errorf("secret deletion err: %w", err)
		}

		return secret, err
	}

	return model.Secret{}, nil
}
