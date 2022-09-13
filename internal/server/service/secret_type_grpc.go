package service

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/server/storage"
	"github.com/sergalkin/gophkeeper/pkg/jwt"
)

type secretTypeGrpc struct {
	pb.UnimplementedSecretServer

	storage storage.SecretTypeServerStorage
	jwt     jwt.Manager
}

func NewSecretTypeGrpc(s storage.SecretTypeServerStorage, j jwt.Manager) *secretTypeGrpc {
	return &secretTypeGrpc{storage: s, jwt: j}
}

func (s *secretTypeGrpc) RegisterService(r grpc.ServiceRegistrar) {
	pb.RegisterSecretServer(r, s)
}

func (s *secretTypeGrpc) GetSecretTypesList(
	ctx context.Context, in *pb.SecretTypesListRequest,
) (*pb.SecretTypesListResponse, error) {
	list, err := s.storage.GetSecretTypes(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &pb.SecretTypesListResponse{}
	for _, secret := range list {
		resp.Secrets = append(resp.Secrets, &pb.SecretType{
			Id:    uint32(secret.ID),
			Title: secret.Title,
		})
	}

	return resp, nil
}
