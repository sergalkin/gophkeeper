package service

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/server/storage"
)

type secretTypeGrpc struct {
	pb.UnimplementedSecretTypeServer

	storage storage.SecretTypeServerStorage
}

// NewSecretTypeGrpc - creates new secret type grpc service.
func NewSecretTypeGrpc(s storage.SecretTypeServerStorage) *secretTypeGrpc {
	return &secretTypeGrpc{storage: s}
}

// RegisterService - registers service via grpc server.
func (s *secretTypeGrpc) RegisterService(r grpc.ServiceRegistrar) {
	pb.RegisterSecretTypeServer(r, s)
}

// GetSecretTypesList - returns list of secret types.
//
// Can be accessed only by authorized users.
func (s *secretTypeGrpc) GetSecretTypesList(
	ctx context.Context, in *pb.SecretTypesListRequest,
) (*pb.SecretTypesListResponse, error) {
	list, err := s.storage.GetSecretTypes(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	resp := &pb.SecretTypesListResponse{}
	for _, secret := range list {
		resp.Secrets = append(resp.Secrets, &pb.Type{
			Id:    uint32(secret.ID),
			Title: secret.Title,
		})
	}

	return resp, nil
}
