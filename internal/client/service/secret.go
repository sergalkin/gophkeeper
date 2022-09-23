package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/client/model"
	"github.com/sergalkin/gophkeeper/internal/client/storage"
	"github.com/sergalkin/gophkeeper/pkg/crypt"
)

type SecretClientService struct {
	glCtx   *model.GlobalContext
	client  pb.SecretClient
	storage storage.Memorier
	crypt   crypt.Crypter
	syncer  storage.Syncer
}

// NewSecretClientService - creates new SecretClientService.
func NewSecretClientService(
	glCtx *model.GlobalContext, client pb.SecretClient, st storage.Memorier, cr crypt.Crypter, sr storage.Syncer,
) *SecretClientService {
	return &SecretClientService{
		glCtx:   glCtx,
		client:  client,
		storage: st,
		crypt:   cr,
		syncer:  sr,
	}
}

// GetListOfSecretes - attempts to return list of secrets from memory, if nothing is found then makes gRPC request
// to server.
func (s *SecretClientService) GetListOfSecretes(id int) ([]*pb.SecretList, error) {
	var list []*pb.SecretList
	list = s.storage.GetSecretList(id)
	if len(list) > 0 {
		return list, nil
	}

	result, err := s.client.GetListOfSecretsByType(s.glCtx.Ctx, &pb.GetListOfSecretsByTypeRequest{TypeId: uint32(id)})
	if err != nil {
		return nil, err
	}

	return result.SecretLists, nil
}

// GetBinarySecret - get binary data from server and stores it into file.
func (s *SecretClientService) GetBinarySecret(id int, location string) error {
	res, err := s.client.GetSecret(s.glCtx.Ctx, &pb.GetSecretRequest{Id: int32(id)})
	if err != nil {
		return err
	}

	if res.Type != 3 {
		return errors.New("this method only works with binary data, please appropriate method next time")
	}

	decoded, errDecode := s.crypt.Decode(string(res.Content))
	if errDecode != nil {
		return errDecode
	}

	f, openErr := os.OpenFile(location, os.O_CREATE|os.O_WRONLY, 0644)
	if openErr != nil {
		return openErr
	}
	defer f.Close()

	_, wrError := f.Write([]byte(decoded))
	if wrError != nil {
		return wrError
	}

	fmt.Println("data written to file")

	return nil
}

// GetSecret - attempts to get secret from memory by its id, if nothing is found then makes gRPC request to server.
func (s *SecretClientService) GetSecret(id int) error {
	data, ok := s.storage.FindInStorage(id)

	if ok {
		fmt.Printf("Content:%+v\n", data)
		return nil
	}

	result, err := s.client.GetSecret(s.glCtx.Ctx, &pb.GetSecretRequest{Id: int32(id)})
	if err != nil {
		return err
	}

	if result.Type == 3 {
		return errors.New("to get binary data, pleas use proper method")
	}

	decoded, errDecode := s.crypt.Decode(string(result.Content))
	if errDecode != nil {
		return errDecode
	}

	var m interface{}
	switch result.Type {
	case 1:
		m = model.LoginPassSecret{}
	case 2:
		m = model.TextSecret{}
	case 4:
		m = model.CardSecret{}
	}

	errUnmarshal := json.Unmarshal([]byte(decoded), &m)
	if errUnmarshal != nil {
		return errUnmarshal
	}

	fmt.Printf("Content:%+v\n", m)

	return nil
}

// CreateSecret - creates new secret on the server and then makes re-sync memory storage.
func (s *SecretClientService) CreateSecret(title string, recordType int, content string) error {
	contentT := []byte(s.crypt.Encode(content))

	result, err := s.client.CreateSecret(s.glCtx.Ctx, &pb.CreateSecretRequest{
		Title:   title,
		Type:    uint32(recordType),
		Content: contentT,
	})

	if err != nil {
		return err
	}

	fmt.Println("created new secret with ID:", result.Id)

	s.syncer.SyncAll()

	return nil
}

// DeleteSecret - deletes a secrete from server and then makes re-sync memory storage.
func (s *SecretClientService) DeleteSecret(id int) error {
	_, err := s.client.DeleteSecret(s.glCtx.Ctx, &pb.DeleteSecretRequest{Id: uint32(id)})
	if err != nil {
		return err
	}

	fmt.Println("successfully deleted secret")

	s.storage.ResetStorage()
	s.syncer.SyncAll()

	return nil
}

// EditSecret - edits secret on the server and then makes re-sync memory storage.
func (s *SecretClientService) EditSecret(id int, title string, recordType int, content string) error {
	contentT := []byte(s.crypt.Encode(content))

	_, err := s.client.EditSecret(
		s.glCtx.Ctx, &pb.EditSecretRequest{Id: uint32(id), Title: title, Type: uint32(recordType), Content: contentT},
	)
	if err != nil {
		return err
	}

	fmt.Println("successfully edited secret")

	s.syncer.SyncAll()

	return nil
}
