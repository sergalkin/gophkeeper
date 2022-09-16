package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/server/middleware/auth"
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

// CreateSecret - stores a provided secret via initialised storage.
func (s *secretGrpc) CreateSecret(ctx context.Context, in *pb.CreateSecretRequest) (*pb.CreateSecretResponse, error) {
	tok := ctx.Value(auth.JwtTokenCtx{}).(string)

	secret := model.Secret{
		UserID:    uuid.MustParse(tok),
		TypeID:    int(in.Type),
		Title:     in.Title,
		Content:   in.Content,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	m, err := s.storage.CreateSecret(ctx, secret)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var deletedAt pb.NullableDeletedAt
	if m.DeletedAt != nil {
		deletedAt = pb.NullableDeletedAt{Kind: &pb.NullableDeletedAt_Data{Data: timestamppb.New(*m.DeletedAt)}}
	} else {
		deletedAt = pb.NullableDeletedAt{Kind: nil}
	}

	return &pb.CreateSecretResponse{
		Id:        uint32(m.ID),
		Title:     m.Title,
		Type:      uint32(m.TypeID),
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
		DeletedAt: &deletedAt,
	}, nil
}

// GetSecret - returns stored secret from storage.
func (s *secretGrpc) GetSecret(ctx context.Context, in *pb.GetSecretRequest) (*pb.GetSecretResponse, error) {
	tok := ctx.Value(auth.JwtTokenCtx{}).(string)

	secret := model.Secret{
		UserID: uuid.MustParse(tok),
		ID:     int(in.Id),
	}

	m, err := s.storage.GetSecret(ctx, secret)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var deletedAt pb.NullableDeletedAt
	if m.DeletedAt != nil {
		deletedAt = pb.NullableDeletedAt{Kind: &pb.NullableDeletedAt_Data{Data: timestamppb.New(*m.DeletedAt)}}
	} else {
		deletedAt = pb.NullableDeletedAt{Kind: nil}
	}

	return &pb.GetSecretResponse{
		Id:        uint32(m.ID),
		Title:     m.Title,
		Type:      uint32(m.TypeID),
		Content:   m.Content,
		CreatedAt: timestamppb.New(m.CreatedAt),
		UpdatedAt: timestamppb.New(m.UpdatedAt),
		DeletedAt: &deletedAt,
	}, nil

}
