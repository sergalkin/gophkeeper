package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/sergalkin/gophkeeper/internal/server/model"
	"github.com/sergalkin/gophkeeper/internal/server/storage"
	"github.com/sergalkin/gophkeeper/pkg/apperr"
)

var _ storage.SecretTypeServerStorage = (*secretTypePostgresStorage)(nil)

type secretTypePostgresStorage struct {
	conn *pgx.Conn
}

const (
	GetSecretTypeList = `SELECT id, title FROM secret_types`
)

func NewPostgresSecretTypeStorage(c *pgx.Conn) *secretTypePostgresStorage {
	return &secretTypePostgresStorage{conn: c}
}

func (s *secretTypePostgresStorage) GetSecretTypes(ctx context.Context) ([]model.SecretType, error) {
	var list []model.SecretType

	rows, err := s.conn.Query(ctx, GetSecretTypeList)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgerrcode.IsNoData(pgErr.Code) {
				return list, apperr.ErrNotFound
			}
		}

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
