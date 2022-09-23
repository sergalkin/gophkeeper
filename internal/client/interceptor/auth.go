package interceptor

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

// AuthInterceptor is a client interceptor for authentication
type AuthInterceptor struct {
	protectedMethods map[string]bool
}

// NewAuthInterceptor - returns an Auth interceptor
func NewAuthInterceptor(prMethods map[string]bool) *AuthInterceptor {
	return &AuthInterceptor{protectedMethods: prMethods}
}

// Unary returns a client interceptor to authenticate unary RPC
func (a *AuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		if a.protectedMethods[method] {
			var token []string
			md, ok := metadata.FromOutgoingContext(ctx)
			if ok {
				token = md.Get("authorization")
			}

			if len(token) == 0 {
				return errors.New("you have to be authorized via login first")
			}

			return invoker(ctx, method, req, reply, cc, opts...)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
