package service

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/server/model"
	"github.com/sergalkin/gophkeeper/internal/server/storage"
)

type secretGrpc struct {
	pb.UnimplementedSecretServer

	storage storage.SecretServerStorage
	//jwtManager jwt.Manager
}

// NewSecretGrpc - creates new secret grpc service.
func NewSecretGrpc(s storage.SecretServerStorage) *secretGrpc {
	return &secretGrpc{
		storage: s,
		//jwtManager: m,
	}
}

// RegisterService - registers service via grpc server.
func (s *secretGrpc) RegisterService(r grpc.ServiceRegistrar) {
	pb.RegisterSecretServer(r, s)
}

func (s *secretGrpc) CreateSecret(ctx context.Context, in *pb.CreateSecretRequest) (*pb.CreateSecretResponse, error) {
	m, err := s.storage.CreateSecret(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.CreateSecretResponse{
		Title: m.Title,
		Type:  uint32(m.ID),
	}, nil
}

func (s *secretGrpc) GetSecret(ctx context.Context, in *pb.GetSecretRequest) (*pb.GetSecretResponse, error) {
	m, err := s.storage.GetSecret(ctx, model.Secret{})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.GetSecretResponse{
		Title:   m.Title,
		Type:    uint32(m.TypeID),
		Content: m.Content,
	}, nil

}
