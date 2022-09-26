package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/server/middleware/auth"
	"github.com/sergalkin/gophkeeper/internal/server/model"
	"github.com/sergalkin/gophkeeper/internal/server/storage"
	"github.com/sergalkin/gophkeeper/pkg/apperr"
)

type SecretGrpc struct {
	pb.UnimplementedSecretServer

	storage storage.SecretServerStorage
}

// NewSecretGrpc - creates new secret grpc service.
func NewSecretGrpc(s storage.SecretServerStorage) *SecretGrpc {
	return &SecretGrpc{
		storage: s,
	}
}

// RegisterService - registers service via grpc server.
func (s *SecretGrpc) RegisterService(r grpc.ServiceRegistrar) {
	pb.RegisterSecretServer(r, s)
}

// CreateSecret - stores a provided secret via initialised storage.
func (s *SecretGrpc) CreateSecret(ctx context.Context, in *pb.CreateSecretRequest) (*pb.CreateSecretResponse, error) {
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
func (s *SecretGrpc) GetSecret(ctx context.Context, in *pb.GetSecretRequest) (*pb.GetSecretResponse, error) {
	tok := ctx.Value(auth.JwtTokenCtx{}).(string)

	secret := model.Secret{
		UserID: uuid.MustParse(tok),
		ID:     int(in.Id),
	}

	m, err := s.storage.GetSecret(ctx, secret)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, err.Error())
		}

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

// DeleteSecret - deletes a user secret from the storage by provide id and userId from ctx.
func (s *SecretGrpc) DeleteSecret(ctx context.Context, in *pb.DeleteSecretRequest) (*pb.DeleteSecretResponse, error) {
	tok := ctx.Value(auth.JwtTokenCtx{}).(string)

	secret := model.Secret{
		UserID: uuid.MustParse(tok),
		ID:     int(in.Id),
	}

	_, err := s.storage.DeleteSecret(ctx, secret)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteSecretResponse{}, nil
}

// EditSecret - edits secret in storage.
func (s *SecretGrpc) EditSecret(ctx context.Context, in *pb.EditSecretRequest) (*pb.EditSecretResponse, error) {
	token := ctx.Value(auth.JwtTokenCtx{}).(string)
	secret := model.Secret{
		ID:        int(in.Id),
		UserID:    uuid.MustParse(token),
		Title:     in.Title,
		TypeID:    int(in.Type),
		Content:   in.Content,
		UpdatedAt: in.UpdatedAt.AsTime(),
	}

	updatedSecret, err := s.storage.EditSecret(ctx, secret, in.IsForce)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		if errors.Is(err, apperr.ErrUpdatedAtDoesntMatch) {
			return nil, status.Error(codes.FailedPrecondition, err.Error())
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	var deletedAt pb.NullableDeletedAt
	if updatedSecret.DeletedAt != nil {
		deletedAt = pb.NullableDeletedAt{Kind: &pb.NullableDeletedAt_Data{
			Data: timestamppb.New(*updatedSecret.DeletedAt)},
		}
	} else {
		deletedAt = pb.NullableDeletedAt{Kind: nil}
	}

	return &pb.EditSecretResponse{
		Id:        uint32(updatedSecret.ID),
		Title:     updatedSecret.Title,
		Type:      uint32(updatedSecret.TypeID),
		CreatedAt: timestamppb.New(updatedSecret.CreatedAt),
		UpdatedAt: timestamppb.New(updatedSecret.UpdatedAt),
		DeletedAt: &deletedAt,
	}, nil
}

// GetListOfSecretsByType - returns list of model.Secret.
func (s *SecretGrpc) GetListOfSecretsByType(
	ctx context.Context, in *pb.GetListOfSecretsByTypeRequest,
) (*pb.GetListOfSecretsByTypeResponse, error) {
	token := ctx.Value(auth.JwtTokenCtx{}).(string)

	userId, errParse := uuid.Parse(token)
	if errParse != nil {
		return nil, status.Error(codes.Internal, errParse.Error())
	}

	user := model.User{ID: &userId}
	secretType := model.SecretType{ID: uint(in.TypeId)}

	secrets, err := s.storage.GetListOfSecretByType(ctx, secretType, user)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var castedSecrets []*pb.SecretList

	for _, val := range secrets {
		var deletedAt pb.NullableDeletedAt
		if val.DeletedAt != nil {
			deletedAt = pb.NullableDeletedAt{Kind: &pb.NullableDeletedAt_Data{Data: timestamppb.New(*val.DeletedAt)}}
		} else {
			deletedAt = pb.NullableDeletedAt{Kind: nil}
		}

		castedSecrets = append(castedSecrets, &pb.SecretList{
			Id:        uint32(val.ID),
			UserId:    val.UserID.String(),
			TypeId:    uint32(val.TypeID),
			Title:     val.Title,
			Content:   val.Content,
			CreatedAt: timestamppb.New(val.CreatedAt),
			UpdatedAt: timestamppb.New(val.UpdatedAt),
			DeletedAt: &deletedAt,
		})
	}

	return &pb.GetListOfSecretsByTypeResponse{SecretLists: castedSecrets}, nil
}
