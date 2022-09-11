package storage

import (
	"context"

	"github.com/sergalkin/gophkeeper/internal/server/model"
)

type UserServerStorage interface {
	Create(ctx context.Context, user model.User) (model.User, error)
	GetByLoginAndPassword(ctx context.Context, login, password string) (user *model.User, err error)
}
