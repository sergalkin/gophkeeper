//go:generate mockgen -source=./interfaces.go -destination=./mock/storage.go -package=storagemock
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
	// DeleteUser - deletes a user from storage.
	DeleteUser(ctx context.Context, user model.User) (model.User, error)
}

type SecretTypeServerStorage interface {
	// GetSecretTypes - returns list of model.SecretType from storage.
	GetSecretTypes(ctx context.Context) ([]model.SecretType, error)
}

type SecretServerStorage interface {
	// CreateSecret - creates new model.Secret in storage.
	CreateSecret(ctx context.Context, secret model.Secret) (model.Secret, error)
	// GetSecret - gets a model.Secret from storage.
	GetSecret(ctx context.Context, secret model.Secret) (model.Secret, error)
	// DeleteSecret - deletes a model.Secret from storage.
	DeleteSecret(ctx context.Context, secret model.Secret) (model.Secret, error)
}
