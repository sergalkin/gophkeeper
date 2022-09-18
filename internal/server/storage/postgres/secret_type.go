package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v4"

	"github.com/sergalkin/gophkeeper/internal/server/model"
	"github.com/sergalkin/gophkeeper/internal/server/storage"
)

var _ storage.SecretTypeServerStorage = (*SecretTypePostgresStorage)(nil)

type SecretTypePostgresStorage struct {
	conn *pgx.Conn
}

const (
	GetSecretTypeList = `SELECT id, title FROM secret_types`
)

// NewPostgresSecretTypeStorage - creates a postgres storage for secret types.
func NewPostgresSecretTypeStorage(c *pgx.Conn) *SecretTypePostgresStorage {
	return &SecretTypePostgresStorage{conn: c}
}

// GetSecretTypes - returns list of all available secret types.
func (s *SecretTypePostgresStorage) GetSecretTypes(ctx context.Context) ([]model.SecretType, error) {
	var list []model.SecretType

	rows, err := s.conn.Query(ctx, GetSecretTypeList)
	if err != nil {
		return list, fmt.Errorf("error in getting secret types list: %w", err)
	}

	for rows.Next() {
		m := model.SecretType{}
		if err = rows.Scan(&m.ID, &m.Title); err != nil {
			return list, fmt.Errorf("error scanning secret types: %w", err)
		}

		list = append(list, m)
	}

	return list, nil
}
