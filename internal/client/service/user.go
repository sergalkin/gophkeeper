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

func NewUserClientService(glCtx *model.GlobalContext, client pb.UserClient) *UserClientService {
	return &UserClientService{
		glCtx:  glCtx,
		client: client,
	}
}

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

func (u *UserClientService) Register(user model.User) error {
	result, err := u.client.Register(u.glCtx.Ctx, &pb.RegisterRequest{Login: user.Login, Password: user.Password})
	if err != nil {
		return err
	}

	u.glCtx.Ctx = metadata.AppendToOutgoingContext(u.glCtx.Ctx, "authorization", "Bearer "+result.Token)

	return nil
}

func (u *UserClientService) Delete() error {
	_, err := u.client.Delete(u.glCtx.Ctx, &pb.DeleteRequest{})
	if err != nil {
		return err
	}

	u.glCtx.Ctx = metadata.NewOutgoingContext(u.glCtx.Ctx, metadata.MD{})

	return nil
}

func (u *UserClientService) Logout() {
	u.glCtx.Ctx = metadata.NewOutgoingContext(u.glCtx.Ctx, metadata.MD{})
}
