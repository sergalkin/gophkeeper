package service

import (
	"google.golang.org/grpc/metadata"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/client/model"
)

type UserClientService struct {
	glCtx  *model.GlobalContext
	client pb.UserClient
}

// NewUserClientService - creates new UserClientService.
func NewUserClientService(glCtx *model.GlobalContext, client pb.UserClient) *UserClientService {
	return &UserClientService{
		glCtx:  glCtx,
		client: client,
	}
}

// Login - authorizes a user by provided login and password. On successful authorization adds authorization token
// to metadata in global shared context.
func (u *UserClientService) Login(user model.User) error {
	request := &pb.LoginRequest{
		Login:    user.Login,
		Password: user.Password,
	}

	result, err := u.client.Login(u.glCtx.Ctx, request)
	if err != nil {
		return err
	}

	u.glCtx.Ctx = metadata.AppendToOutgoingContext(u.glCtx.Ctx, "authorization", "Bearer "+result.Token)

	return nil
}

// Register - creates a new user on server. On successful creation adds authorization token to metadata in global
// shared context.
func (u *UserClientService) Register(user model.User) error {
	result, err := u.client.Register(u.glCtx.Ctx, &pb.RegisterRequest{Login: user.Login, Password: user.Password})
	if err != nil {
		return err
	}

	u.glCtx.Ctx = metadata.AppendToOutgoingContext(u.glCtx.Ctx, "authorization", "Bearer "+result.Token)

	return nil
}

// Delete - deletes a user from server. On successful deletion, removes authorization token from metadata in global
// shared context.
func (u *UserClientService) Delete() error {
	_, err := u.client.Delete(u.glCtx.Ctx, &pb.DeleteRequest{})
	if err != nil {
		return err
	}

	u.glCtx.Ctx = metadata.NewOutgoingContext(u.glCtx.Ctx, metadata.MD{})

	return nil
}

// Logout -  removes authorization token from metadata in global shared context.
func (u *UserClientService) Logout() {
	u.glCtx.Ctx = metadata.NewOutgoingContext(u.glCtx.Ctx, metadata.MD{})
}
