package service

import (
	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/client/model"
)

type SecretTypeClientService struct {
	glCtx  *model.GlobalContext
	client pb.SecretTypeClient
}

// NewSecretTypeClientService - creates new SecretTypeClientService.
func NewSecretTypeClientService(glCtx *model.GlobalContext, client pb.SecretTypeClient) *SecretTypeClientService {
	return &SecretTypeClientService{
		glCtx:  glCtx,
		client: client,
	}
}

// List - returns secret types list from server.
func (s *SecretTypeClientService) List() (*pb.SecretTypesListResponse, error) {
	request := &pb.SecretTypesListRequest{}

	result, err := s.client.GetSecretTypesList(s.glCtx.Ctx, request)
	if err != nil {
		return nil, err
	}

	return result, nil
}
