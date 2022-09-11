package postgres

import (
	"context"

	"github.com/jackc/pgx/v4"

	"github.com/sergalkin/gophkeeper/internal/server/model"
	"github.com/sergalkin/gophkeeper/internal/server/storage"
)

var _ storage.UserServerStorage = (*userPostgresStorage)(nil)

type userPostgresStorage struct {
	conn *pgx.Conn
}

func NewPostgresUserStorage(c *pgx.Conn) *userPostgresStorage {
	return &userPostgresStorage{conn: c}
}

func (u userPostgresStorage) Create(ctx context.Context, user model.User) (model.User, error) {
	// TODO implement me
	panic("implement me")
}

func (u userPostgresStorage) GetByLoginAndPassword(ctx context.Context, login, password string) (*model.User, error) {
	// TODO implement me
	panic("implement me")
}
