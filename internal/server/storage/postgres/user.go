package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v4"

	"github.com/sergalkin/gophkeeper/internal/server/model"
	"github.com/sergalkin/gophkeeper/internal/server/storage"
	"github.com/sergalkin/gophkeeper/pkg/apperr"
)

var _ storage.UserServerStorage = (*userPostgresStorage)(nil)

type userPostgresStorage struct {
	conn *pgx.Conn
}

const (
	CreateUser = `INSERT INTO users (login, password) VALUES ($1, crypt($2, gen_salt('bf'))) returning id`
	GetUserId  = `SELECT id FROM users WHERE login = $1 AND password = crypt($2, password)`
)

// NewPostgresUserStorage - Creates userPostgresStorage instance.
func NewPostgresUserStorage(c *pgx.Conn) *userPostgresStorage {
	return &userPostgresStorage{conn: c}
}

// Create - creates a user record in DB with data provided from model.User, then returns model.User populated with
// user id from database.
func (u userPostgresStorage) Create(ctx context.Context, user model.User) (model.User, error) {
	ctxWithTimeOut, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := u.conn.QueryRow(ctxWithTimeOut, CreateUser, user.Login, user.Password).Scan(&user.ID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
				return user, apperr.ErrConflict
			}
		}

		return user, fmt.Errorf("user insertion err: %w", err)
	}

	return user, nil
}

// GetByLoginAndPassword - searches DB with provided model.User, if record is found, then populates model.User with
// user id from database.
func (u userPostgresStorage) GetByLoginAndPassword(ctx context.Context, user model.User) (model.User, error) {
	ctxWithTimeOut, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	err := u.conn.QueryRow(ctxWithTimeOut, GetUserId, user.Login, user.Password).Scan(&user.ID)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			if pgerrcode.IsNoData(pgErr.Code) {
				return user, apperr.ErrNotFound
			}
		}

		return user, fmt.Errorf("user insertion err: %w", err)
	}

	return user, nil
}
