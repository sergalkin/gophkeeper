package storage

import (
	"encoding/json"
	"fmt"

	pb "github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/client/model"
	"github.com/sergalkin/gophkeeper/pkg/crypt"
)

type Syncer interface {
	SyncAll()
	SyncPassLoginData() error
	SyncCardData() error
	SyncTextData() error
}

type Sync struct {
	storage      Memorier
	secretClient pb.SecretClient
	glCtx        *model.GlobalContext
	cr           crypt.Crypter
}

// NewSync - creates new Sync.
func NewSync(s Memorier, sc pb.SecretClient, ctx *model.GlobalContext, cr crypt.Crypter) *Sync {
	return &Sync{storage: s, secretClient: sc, glCtx: ctx, cr: cr}
}

// SyncAll - runs SyncTextData, SyncPassLoginData, SyncCardData under the hood.
func (s *Sync) SyncAll() {
	if err := s.SyncTextData(); err != nil {
		fmt.Println(err)
	}

	if err := s.SyncPassLoginData(); err != nil {
		fmt.Println(err)
	}

	if err := s.SyncCardData(); err != nil {
		fmt.Println(err)
	}
}

// SyncTextData - makes gRPC request to server and on success sets acquired records to MemoryStorage.TextSecrets.
func (s *Sync) SyncTextData() error {
	texts, err := s.secretClient.GetListOfSecretsByType(s.glCtx.Ctx, &pb.GetListOfSecretsByTypeRequest{TypeId: 2})
	if err != nil {
		panic(err)
	}

	var list []model.TextSecret
	for _, text := range texts.SecretLists {
		id := int(text.Id)
		m := model.TextSecret{
			Id:         id,
			Title:      text.Title,
			RecordType: 2,
		}

		decoded, errDecode := s.cr.Decode(string(text.Content))
		if errDecode != nil {
			return errDecode
		}

		errUnmarshal := json.Unmarshal([]byte(decoded), &m)
		if errUnmarshal != nil {
			return errUnmarshal
		}
		m.Id = id

		list = append(list, m)
	}

	s.storage.SetTextSecrets(list)

	return nil
}

// SyncCardData - makes gRPC request to server and on success sets acquired records to MemoryStorage.CardSecrets.
func (s *Sync) SyncCardData() error {
	cards, err := s.secretClient.GetListOfSecretsByType(s.glCtx.Ctx, &pb.GetListOfSecretsByTypeRequest{TypeId: 4})
	if err != nil {
		panic(err)
	}

	var list []model.CardSecret
	for _, card := range cards.SecretLists {
		id := int(card.Id)
		m := model.CardSecret{
			Id:         int(card.Id),
			Title:      card.Title,
			RecordType: 4,
		}

		decoded, errDecode := s.cr.Decode(string(card.Content))
		if errDecode != nil {
			return errDecode
		}

		errUnmarshal := json.Unmarshal([]byte(decoded), &m)
		if errUnmarshal != nil {
			return errUnmarshal
		}
		m.Id = id

		list = append(list, m)
	}

	s.storage.SetCardSecrets(list)

	return nil
}

// SyncPassLoginData - makes gRPC request to server and on success sets acquired records to
// MemoryStorage.LoginPassSecrets.
func (s *Sync) SyncPassLoginData() error {
	lists, err := s.secretClient.GetListOfSecretsByType(s.glCtx.Ctx, &pb.GetListOfSecretsByTypeRequest{TypeId: 1})
	if err != nil {
		panic(err)
	}

	var list []model.LoginPassSecret
	for _, sList := range lists.SecretLists {
		id := int(sList.Id)
		m := model.LoginPassSecret{
			Id:         int(sList.Id),
			Title:      sList.Title,
			RecordType: int(sList.TypeId),
		}

		decoded, errDecode := s.cr.Decode(string(sList.Content))
		if errDecode != nil {
			return errDecode
		}

		errUnmarshal := json.Unmarshal([]byte(decoded), &m)
		if errUnmarshal != nil {
			return errUnmarshal
		}
		m.Id = id

		list = append(list, m)
	}

	s.storage.SetLoginPassSecrets(list)

	return nil
}
