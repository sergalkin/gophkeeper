package storage

import (
	"context"

	"github.com/sergalkin/gophkeeper/internal/server/model"
)

type UserServerStorage interface {
	// Create - create a new model.User in storage.
	Create(ctx context.Context, user model.User) (model.User, error)
	// GetByLoginAndPassword - returns model.User from storage.
	GetByLoginAndPassword(ctx context.Context, user model.User) (model.User, error)
}

type SecretTypeServerStorage interface {
	// GetSecretTypes - returns list of model.SecretType from storage.
	GetSecretTypes(ctx context.Context) ([]model.SecretType, error)
}

type SecretServerStorage interface {
}
