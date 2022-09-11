package service

import (
	"context"
	"errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/server/model"
	"github.com/sergalkin/gophkeeper/internal/server/storage"
	"github.com/sergalkin/gophkeeper/pkg/apperr"
)

type userGrpc struct {
	pb.UnimplementedUserServer

	storage storage.UserServerStorage
}

func NewUserGrpc(s storage.UserServerStorage) *userGrpc {
	return &userGrpc{
		storage: s,
	}
}

func (u *userGrpc) RegisterService(r grpc.ServiceRegistrar) {
	pb.RegisterUserServer(r, u)
}

func (u *userGrpc) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	userModel, err := u.storage.Create(ctx, model.User{Login: in.Login, Password: in.Password})

	if errors.Is(err, apperr.ErrConflict) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{Token: userModel.Identity()}, nil
}
