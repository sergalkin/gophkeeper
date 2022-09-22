package service

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/server/middleware/auth"
	"github.com/sergalkin/gophkeeper/internal/server/model"
	storagemock "github.com/sergalkin/gophkeeper/internal/server/storage/mock"
	"github.com/sergalkin/gophkeeper/pkg/apperr"
	cryptmock "github.com/sergalkin/gophkeeper/pkg/crypt/mock"
	jwtmock "github.com/sergalkin/gophkeeper/pkg/jwt/mock"
)

func Test_userGrpc_RegisterService(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	userMock := storagemock.NewMockUserServerStorage(ctl)
	jwtMock := jwtmock.NewMockManager(ctl)
	cryptMock := cryptmock.NewMockCrypter(ctl)

	tests := []struct {
		name string
	}{
		{
			name: "Registrar can be called without errors",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewUserGrpc(userMock, jwtMock, cryptMock)

			server := grpc.NewServer()

			s.RegisterService(server)
		})
	}
}

func Test_userGrpc_Register(t *testing.T) {
	uid := uuid.New()
	ctx := context.WithValue(context.Background(), auth.JwtTokenCtx{}, uid)

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	client, done := userTestClient(t, ctl, uid)
	defer close(done)

	_, err := client.Register(ctx, &pb.RegisterRequest{
		Login:    "RegConflict",
		Password: "pass",
	})
	assert.Error(t, err)

	_, err = client.Register(ctx, &pb.RegisterRequest{
		Login:    "RegOk",
		Password: "pass",
	})
	assert.NoError(t, err)
}

func Test_userGrpc_Login(t *testing.T) {
	uid := uuid.New()
	ctx := context.WithValue(context.Background(), auth.JwtTokenCtx{}, uid)

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	client, done := userTestClient(t, ctl, uid)
	defer close(done)

	_, err := client.Login(ctx, &pb.LoginRequest{
		Login:    "loginOk",
		Password: "pass",
	})
	assert.NoError(t, err)

	_, err = client.Login(ctx, &pb.LoginRequest{
		Login:    "logineErr",
		Password: "pass",
	})
	assert.Error(t, err)
}

func Test_userGrpc_Delete(t *testing.T) {
	uid := uuid.New()

	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer token")

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	client, done := userTestClient(t, ctl, uid)
	defer close(done)

	_, err := client.Delete(ctx, &pb.DeleteRequest{})
	assert.NoError(t, err)
}

func userTestClient(t *testing.T, ctl *gomock.Controller, uid uuid.UUID) (pb.UserClient, chan<- struct{}) {
	done := make(chan struct{})

	userStorageMock := storagemock.NewMockUserServerStorage(ctl)
	userStorageMock.
		EXPECT().
		GetByLoginAndPassword(gomock.Any(), gomock.Eq(model.User{Login: "loginOk", Password: "pass"})).
		AnyTimes().
		Return(model.User{ID: &uid, Login: "test", Password: "pass"}, nil)

	userStorageMock.
		EXPECT().
		GetByLoginAndPassword(gomock.Any(), gomock.Eq(model.User{Login: "logineErr", Password: "pass"})).
		AnyTimes().
		Return(model.User{Login: "test", Password: "pass"}, errors.New("user login error"))

	userStorageMock.
		EXPECT().
		Create(gomock.Any(), gomock.Eq(model.User{Login: "RegConflict", Password: "pass"})).
		AnyTimes().
		Return(model.User{Login: "test", Password: "pass"}, apperr.ErrConflict)

	userStorageMock.
		EXPECT().
		Create(gomock.Any(), gomock.Eq(model.User{Login: "RegOk", Password: "pass"})).
		AnyTimes().
		Return(model.User{ID: &uid, Login: "test1", Password: "pass"}, nil)

	userStorageMock.
		EXPECT().
		DeleteUser(
			gomock.Any(),
			gomock.Any(),
		).
		AnyTimes().
		Return(model.User{}, nil)

	jwtM := jwtmock.NewMockManager(ctl)
	jwtM.EXPECT().Issue(uid.String()).AnyTimes().Return("token", nil)
	jwtM.EXPECT().Decode(gomock.Any()).AnyTimes().Return(uid.String(), nil)

	cryptM := cryptmock.NewMockCrypter(ctl)
	cryptM.EXPECT().Encode(gomock.Any()).AnyTimes().Return("token")
	cryptM.EXPECT().Decode(gomock.Any()).AnyTimes().Return(uid.String(), nil)

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}

	userRpc := NewUserGrpc(userStorageMock, jwtM, cryptM)

	server := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpcauth.UnaryServerInterceptor(auth.NewJwtMiddleware(jwtM, cryptM).Auth),
			)),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpcauth.StreamServerInterceptor(auth.NewJwtMiddleware(jwtM, cryptM).Auth),
			)),
	)

	pb.RegisterUserServer(server, userRpc)

	go func() {
		if err = server.Serve(l); err != nil && err != grpc.ErrServerStopped {
			panic(err)
		}
	}()

	conn, err := grpc.Dial(l.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}

	go func() {
		for {
			select {
			case <-done:
				server.GracefulStop()
				conn.Close()
			default:
			}
		}
	}()

	client := pb.NewUserClient(conn)

	return client, done
}
