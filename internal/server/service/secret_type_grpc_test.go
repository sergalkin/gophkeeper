package service

import (
	"context"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/server/model"
	storagemock "github.com/sergalkin/gophkeeper/internal/server/storage/mock"
)

func Test_secretTypeGrpc_GetSecretTypesList(t *testing.T) {
	ctx := context.Background()

	ctl := gomock.NewController(t)
	defer ctl.Finish()

	secretTypeMock := storagemock.NewMockSecretTypeServerStorage(ctl)

	secretTypeMock.EXPECT().GetSecretTypes(gomock.Any()).AnyTimes().Return(
		[]model.SecretType{
			{ID: 1, Title: "Test"},
			{ID: 2, Title: "Test 2"},
		},
		nil,
	)

	l, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal(err)
	}

	secretTypeRpc := NewSecretTypeGrpc(secretTypeMock)

	server := grpc.NewServer()
	defer server.GracefulStop()

	pb.RegisterSecretTypeServer(server, secretTypeRpc)

	go func() {
		if err = server.Serve(l); err != nil && err != grpc.ErrServerStopped {
			panic(err)
		}
	}()

	conn, err := grpc.Dial(l.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	client := pb.NewSecretTypeClient(conn)

	resp, err := client.GetSecretTypesList(ctx, &pb.SecretTypesListRequest{})
	assert.NoError(t, err)
	assert.Len(t, resp.Secrets, 2)
}

func Test_secretTypeGrpc_RegisterService(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	secretTypeMock := storagemock.NewMockSecretTypeServerStorage(ctl)

	tests := []struct {
		name string
	}{
		{
			name: "Registrar can be called without errors",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewSecretTypeGrpc(secretTypeMock)

			server := grpc.NewServer()

			s.RegisterService(server)
		})
	}
}
