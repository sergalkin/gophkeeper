package service

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/server/middleware/auth"
	"github.com/sergalkin/gophkeeper/internal/server/model"
	storagemock "github.com/sergalkin/gophkeeper/internal/server/storage/mock"
	cryptmock "github.com/sergalkin/gophkeeper/pkg/crypt/mock"
	jwtmock "github.com/sergalkin/gophkeeper/pkg/jwt/mock"
)

var loc, _ = time.LoadLocation("UTC")
var now = time.Now().In(loc)

func TestSecretGrpc_RegisterService(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	secretMock := storagemock.NewMockSecretServerStorage(ctl)

	tests := []struct {
		name string
	}{
		{
			name: "Registrar can be called without errors",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSecretGrpc(secretMock)

			server := grpc.NewServer()

			s.RegisterService(server)
		})
	}
}

func TestSecretGrpc_CreateSecret(t *testing.T) {
	uid := uuid.New()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer token")

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	client, done := secretTestClient(t, ctl, uid)
	defer close(done)

	_, err := client.CreateSecret(ctx, &pb.CreateSecretRequest{})
	assert.NoError(t, err)
}

func TestSecretGrpc_GetSecret(t *testing.T) {
	uid := uuid.New()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer token")

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	client, done := secretTestClient(t, ctl, uid)
	defer close(done)

	_, err := client.GetSecret(ctx, &pb.GetSecretRequest{Id: 1})
	assert.NoError(t, err)

	_, err = client.GetSecret(ctx, &pb.GetSecretRequest{Id: 0})
	assert.Error(t, err)
}

func TestSecretGrpc_DeleteSecret(t *testing.T) {
	uid := uuid.New()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer token")

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	client, done := secretTestClient(t, ctl, uid)
	defer close(done)

	_, err := client.DeleteSecret(ctx, &pb.DeleteSecretRequest{Id: 1})
	assert.NoError(t, err)

	_, err = client.DeleteSecret(ctx, &pb.DeleteSecretRequest{Id: 0})
	assert.Error(t, err)
}

func TestSecretGrpc_EditSecret(t *testing.T) {
	uid := uuid.New()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer token")

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	client, done := secretTestClient(t, ctl, uid)
	defer close(done)

	_, err := client.EditSecret(ctx, &pb.EditSecretRequest{Id: 1, IsForce: true, UpdatedAt: timestamppb.New(now)})
	assert.NoError(t, err)

	_, err = client.EditSecret(ctx, &pb.EditSecretRequest{Id: 0, IsForce: true, UpdatedAt: timestamppb.New(now)})
	assert.Error(t, err)
}

func TestSecretGrpc_GetListOfSecretsByType(t *testing.T) {
	uid := uuid.New()
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", "Bearer token")

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	client, done := secretTestClient(t, ctl, uid)
	defer close(done)

	res, err := client.GetListOfSecretsByType(ctx, &pb.GetListOfSecretsByTypeRequest{})
	assert.NoError(t, err)
	assert.Len(t, res.SecretLists, 2)
}

func secretTestClient(t *testing.T, ctl *gomock.Controller, uid uuid.UUID) (pb.SecretClient, chan<- struct{}) {
	done := make(chan struct{})

	secretStorageMock := storagemock.NewMockSecretServerStorage(ctl)

	secretStorageMock.EXPECT().CreateSecret(gomock.Any(), gomock.Any()).AnyTimes().Return(model.Secret{}, nil)

	secretStorageMock.EXPECT().
		GetSecret(gomock.Any(), gomock.Eq(model.Secret{ID: 1, UserID: uid})).
		AnyTimes().
		Return(model.Secret{}, nil)
	secretStorageMock.EXPECT().
		GetSecret(gomock.Any(), gomock.Eq(model.Secret{ID: 0, UserID: uid})).
		AnyTimes().
		Return(model.Secret{}, errors.New("test"))

	secretStorageMock.EXPECT().
		DeleteSecret(gomock.Any(), gomock.Eq(model.Secret{ID: 1, UserID: uid})).
		AnyTimes().
		Return(model.Secret{}, nil)
	secretStorageMock.EXPECT().
		DeleteSecret(gomock.Any(), gomock.Eq(model.Secret{ID: 0, UserID: uid})).
		AnyTimes().
		Return(model.Secret{}, errors.New("test"))

	secretStorageMock.EXPECT().
		EditSecret(gomock.Any(), gomock.Eq(model.Secret{ID: 1, UserID: uid, UpdatedAt: now}), gomock.Eq(true)).
		AnyTimes().
		Return(model.Secret{}, nil)
	secretStorageMock.EXPECT().
		EditSecret(gomock.Any(), gomock.Eq(model.Secret{ID: 0, UserID: uid, UpdatedAt: now}), gomock.Eq(true)).
		AnyTimes().
		Return(model.Secret{}, errors.New("test"))

	secretStorageMock.EXPECT().GetListOfSecretByType(gomock.Any(), gomock.Any(), gomock.Any()).AnyTimes().
		Return([]model.Secret{
			{ID: 1},
			{ID: 2},
		}, nil)

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

	secretRpc := NewSecretGrpc(secretStorageMock)

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

	pb.RegisterSecretServer(server, secretRpc)

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

	client := pb.NewSecretClient(conn)

	return client, done
}
