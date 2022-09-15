package service

import (
	"context"
	"errors"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/server/model"
	"github.com/sergalkin/gophkeeper/internal/server/storage"
	"github.com/sergalkin/gophkeeper/pkg/apperr"
	"github.com/sergalkin/gophkeeper/pkg/jwt"
)

type userGrpc struct {
	pb.UnimplementedUserServer

	storage    storage.UserServerStorage
	jwtManager jwt.Manager
}

// NewUserGrpc - creates new user grpc service.
func NewUserGrpc(s storage.UserServerStorage, m jwt.Manager) *userGrpc {
	return &userGrpc{
		storage:    s,
		jwtManager: m,
	}
}

// RegisterService - registers service via grpc server.
func (u *userGrpc) RegisterService(r grpc.ServiceRegistrar) {
	pb.RegisterUserServer(r, u)
}

// Register - registers a new user.
//
//On successful creation returns JwtToken.
func (u *userGrpc) Register(ctx context.Context, in *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	m := model.User{Login: in.Login, Password: in.Password}

	validate := validator.New()
	if errV := validate.Struct(m); errV != nil {
		return nil, status.Error(codes.InvalidArgument, errV.Error())
	}

	userModel, err := u.storage.Create(ctx, m)

	if errors.Is(err, apperr.ErrConflict) {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	token, errToken := u.jwtManager.Issue(userModel.ID.String())
	if errToken != nil {
		return nil, status.Error(codes.Internal, errToken.Error())
	}

	return &pb.RegisterResponse{Token: token}, nil
}

// Login - Will return JwtToken on successful authentication via provided login and password.
func (u *userGrpc) Login(ctx context.Context, in *pb.LoginRequest) (*pb.LoginResponse, error) {
	userModel, err := u.storage.GetByLoginAndPassword(ctx, model.User{Login: in.Login, Password: in.Password})

	if errors.Is(err, apperr.ErrNotFound) {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	token, errToken := u.jwtManager.Issue(userModel.ID.String())
	if errToken != nil {
		return nil, status.Error(codes.Internal, errToken.Error())
	}

	return &pb.LoginResponse{Token: token}, nil
}
