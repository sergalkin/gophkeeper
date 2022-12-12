package auth

import "context"

type Auther interface {
	Auth(ctx context.Context) (context.Context, error)
}
