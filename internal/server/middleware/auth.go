package middleware

import (
	"context"

	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/sergalkin/gophkeeper/pkg/jwt"
)

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

// JwtAuth - a jwt middleware that validates that current requested method can be accessed only by authorized users
// and then checks that bearer token in authorization header is present and is valid jwt token.
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

	newCtx := context.WithValue(ctx, "token", decodedToken)

	return newCtx, nil
}

// isSkippingCurrentRoute - is helper function for checking that current requested method is closed by authorized only
// users.
//
// If requested method is protected by authorization than returns false, otherwise returns true.
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
