package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/sergalkin/gophkeeper/internal/server/model"
	"github.com/sergalkin/gophkeeper/internal/server/storage"
)

var _ storage.SecretServerStorage = (*secretPostgresStorage)(nil)

type secretPostgresStorage struct {
	conn *pgx.Conn
}

func NewSecretPostgresStorage(c *pgx.Conn) *secretPostgresStorage {
	return &secretPostgresStorage{conn: c}
}

func (s *secretPostgresStorage) CreateSecret(ctx context.Context) (model.Secret, error) {
	return model.Secret{}, nil
}

func (s *secretPostgresStorage) GetSecret(ctx context.Context, secret model.Secret) (model.Secret, error) {
	return model.Secret{}, nil
}
