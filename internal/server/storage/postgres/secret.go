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
	"github.com/sergalkin/gophkeeper/pkg/apperr"
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
	DeleteSecret = `delete from secrets where id = $1 and user_id = $2 returning id`
	UpdateSecret = `update secrets 
					set title = $1, content = $2, updated_at = $3
					where id = $4 and user_id = $5
					returning type_id, updated_at, deleted_at
`
	SecretsByType = `select id, user_id, type_id, title, content, created_at, updated_at, deleted_at 
					 from secrets
					 where type_id = $1 and user_id = $2 
`
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
		if !errors.Is(err, pgx.ErrNoRows) {
			return secret, fmt.Errorf("secret getting error: %w", err)
		}

		return secret, fmt.Errorf("error in getting secret from db: %w", err)
	}

	decode, decErr := hex.DecodeString(string(secret.Content))
	if decErr != nil {
		return secret, fmt.Errorf("error in decoding content from db: %w", decErr)
	}

	secret.Content = decode

	return secret, nil
}

// DeleteSecret - deletes a model.Secret from database.
func (s *SecretPostgresStorage) DeleteSecret(ctx context.Context, secret model.Secret) (model.Secret, error) {
	ctxWithTimeOut, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var deletedSecretId *int
	err := s.conn.QueryRow(ctxWithTimeOut, DeleteSecret, secret.ID, secret.UserID).Scan(&deletedSecretId)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return secret, fmt.Errorf("secret deletion err: %w", err)
		}

		return secret, err
	}

	return model.Secret{}, nil
}

// EditSecret - updates a model.Secret in database.
func (s *SecretPostgresStorage) EditSecret(ctx context.Context, secret model.Secret, isForce bool) (model.Secret, error) {
	ctxWithTimeOut, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	ss, _ := s.GetSecret(ctx, secret)

	if ss.UpdatedAt.Unix() != secret.UpdatedAt.Unix() && !isForce {
		return secret, apperr.ErrUpdatedAtDoesntMatch
	}

	err := s.conn.QueryRow(ctxWithTimeOut, UpdateSecret, secret.Title, hex.EncodeToString(secret.Content),
		time.Now(), secret.ID, secret.UserID,
	).Scan(&secret.TypeID, &secret.UpdatedAt, &secret.DeletedAt)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return secret, fmt.Errorf("secret updating error: %w", err)
		}

		return secret, err
	}

	return secret, nil
}

// GetListOfSecretByType - returns a []model.Secret from database by provided type_id via model.SecretType and user_id
// via model.User.
func (s *SecretPostgresStorage) GetListOfSecretByType(
	ctx context.Context, secretType model.SecretType, user model.User,
) ([]model.Secret, error) {
	ctxWithTimeOut, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var secrets []model.Secret

	rows, err := s.conn.Query(ctxWithTimeOut, SecretsByType, secretType.ID, user.ID)
	if err != nil {
		return secrets, fmt.Errorf("getting list of secrets error: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var secret model.Secret

		if scanErr := rows.Scan(
			&secret.ID,
			&secret.UserID,
			&secret.TypeID,
			&secret.Title,
			&secret.Content,
			&secret.CreatedAt,
			&secret.UpdatedAt,
			&secret.DeletedAt,
		); scanErr != nil {
			return secrets, fmt.Errorf("error in scanning gotten row: %w", scanErr)
		}

		secret.Content, err = hex.DecodeString(string(secret.Content))

		if err != nil {
			return secrets, fmt.Errorf("error in decodeing content from row: %w", err)
		}

		secrets = append(secrets, secret)
	}

	return secrets, nil
}
