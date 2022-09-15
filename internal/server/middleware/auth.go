package middleware

import (
	"context"

	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sergalkin/gophkeeper/pkg/jwt"
)

// JwtTokenCtx - a unique type to avoid collisions.
type JwtTokenCtx struct{}

type AuthMiddleware struct {
	jwtManager         jwt.Manager
	unProtectedMethods []string
}

// NewAuthMiddleware - creates AuthMiddleware.
func NewAuthMiddleware(j jwt.Manager) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager:         j,
		unProtectedMethods: []string{"/proto.User/Register", "/proto.User/Login"},
	}
}

// JwtAuth - middleware function for validation user JwtToken.
//
// It extracts and decodes bearer token from context.
//
// On successful decode it attaches decoded token to context with new value with JwtTokenCtx.
// On failure attempt returns codes.Unauthenticated status.
func (a *AuthMiddleware) JwtAuth(ctx context.Context) (context.Context, error) {
	if a.isSkippingCurrentRoute(ctx) {
		return ctx, nil
	}

	token, err := grpcauth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	decodedToken, errDecode := a.jwtManager.Decode(token)
	if errDecode != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", errDecode)
	}

	newCtx := context.WithValue(ctx, JwtTokenCtx{}, decodedToken)

	return newCtx, nil
}

// isSkippingCurrentRoute - helper function for validating that extracted name of currently requested grpc.Method
// from context is in list of unprotected AuthMiddleware a methods.
func (a *AuthMiddleware) isSkippingCurrentRoute(ctx context.Context) bool {
	isSkipping := false

	calledMethod, _ := grpc.Method(ctx)
	for _, method := range a.unProtectedMethods {
		if calledMethod == method {
			isSkipping = true
		}
	}

	return isSkipping
}
