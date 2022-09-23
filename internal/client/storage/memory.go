package storage

import (
	"errors"
	"sync"

	"github.com/sergalkin/gophkeeper/api/proto"
	"github.com/sergalkin/gophkeeper/internal/client/model"
)

type Memorier interface {
	GetLoginPassSecret(id int) (model.LoginPassSecret, bool, error)
	SetLoginPassSecrets([]model.LoginPassSecret)
	GetCardSecret(id int) (model.CardSecret, bool, error)
	SetCardSecrets([]model.CardSecret)
	GetTextSecret(id int) (model.TextSecret, bool, error)
	SetTextSecrets([]model.TextSecret)
	FindInStorage(id int) (interface{}, bool)
	GetSecretList(id int) []*proto.SecretList
	ResetStorage()
}

type MemoryStorage struct {
	mu               sync.RWMutex
	LoginPassSecrets map[int]model.LoginPassSecret
	TextSecrets      map[int]model.TextSecret
	CardSecrets      map[int]model.CardSecret
}

// NewMemoryStorage - creates new MemoryStorage.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		LoginPassSecrets: make(map[int]model.LoginPassSecret, 0),
		TextSecrets:      make(map[int]model.TextSecret, 0),
		CardSecrets:      make(map[int]model.CardSecret, 0),
	}
}

// GetLoginPassSecret - attempts to return model.LoginPassSecret from MemoryStorage.
//
// If record is found, then returns model.LoginPassSecret and true as representation of found record.
//
// If nothing is found, then returns empty model.LoginPassSecret and false as representation of not found record.
func (ms *MemoryStorage) GetLoginPassSecret(id int) (model.LoginPassSecret, bool, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	data, ok := ms.LoginPassSecrets[id]
	if !ok {
		return model.LoginPassSecret{}, ok, errors.New("Login/Pass not found")
	}

	return data, ok, nil
}

// SetLoginPassSecrets - Sets []model.LoginPassSecret to MemoryStorage.
func (ms *MemoryStorage) SetLoginPassSecrets(models []model.LoginPassSecret) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	for _, m := range models {
		ms.LoginPassSecrets[m.Id] = m
	}
}

// ResetStorage - removes all records from MemoryStorage.
func (ms *MemoryStorage) ResetStorage() {
	ms.LoginPassSecrets = make(map[int]model.LoginPassSecret, 0)
	ms.TextSecrets = make(map[int]model.TextSecret, 0)
	ms.CardSecrets = make(map[int]model.CardSecret, 0)
}

// SetCardSecrets - Sets []model.CardSecret to MemoryStorage.
func (ms *MemoryStorage) SetCardSecrets(models []model.CardSecret) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	for _, m := range models {
		ms.CardSecrets[m.Id] = m
	}
}

// GetCardSecret - attempts to return model.CardSecret from MemoryStorage.
//
// If record is found, then returns model.CardSecret and true as representation of found record.
//
// If nothing is found, then returns empty model.CardSecret and false as representation of not found record.
func (ms *MemoryStorage) GetCardSecret(id int) (model.CardSecret, bool, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	data, ok := ms.CardSecrets[id]
	if !ok {
		return model.CardSecret{}, ok, errors.New("card data not found")
	}

	return data, ok, nil
}

// SetTextSecrets - Sets []model.TextSecret to MemoryStorage.
func (ms *MemoryStorage) SetTextSecrets(models []model.TextSecret) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()

	for _, m := range models {
		ms.TextSecrets[m.Id] = m
	}
}

// GetTextSecret - attempts to return model.TextSecret from MemoryStorage.
//
// If record is found, then returns model.TextSecret and true as representation of found record.
//
// If nothing is found, then returns empty model.TextSecret and false as representation of not found record.
func (ms *MemoryStorage) GetTextSecret(id int) (model.TextSecret, bool, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	data, ok := ms.TextSecrets[id]
	if !ok {
		return model.TextSecret{}, ok, errors.New("text not found")
	}

	return data, ok, nil
}

// FindInStorage - attempts to find record in MemoryStorage by provided id.
//
// If record is found in LoginPassSecrets then returns model.LoginPassSecret and true as representation of found record.
//
// If record is found in TextSecrets then returns model.TextSecret and true as representation of found record.
//
// If record is found in CardSecrets then returns model.CardSecret and true as representation of found record.
//
// Otherwise, returns nil and false as representation of not found record.
func (ms *MemoryStorage) FindInStorage(id int) (interface{}, bool) {
	data, ok, _ := ms.GetLoginPassSecret(id)
	if ok {
		return data, true
	}

	dataCard, okCard, _ := ms.GetCardSecret(id)
	if okCard {
		return dataCard, true
	}

	textData, okText, _ := ms.GetTextSecret(id)
	if okText {
		return textData, true
	}

	return nil, false
}

// GetSecretList - goes through all MemoryStorage records and returns []*GetSecretList.
func (ms *MemoryStorage) GetSecretList(id int) []*proto.SecretList {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	var list []*proto.SecretList
	switch id {
	case 1:
		for _, data := range ms.LoginPassSecrets {
			list = append(list, &proto.SecretList{Id: uint32(data.Id), Title: data.Title})
		}
	case 2:
		for _, data := range ms.TextSecrets {
			list = append(list, &proto.SecretList{Id: uint32(data.Id), Title: data.Title})
		}
	case 4:
		for _, data := range ms.CardSecrets {
			list = append(list, &proto.SecretList{Id: uint32(data.Id), Title: data.Title})
		}
	}
	//
	//
	//fmt.Println(data)
	//for _, data := range ms.LoginPassSecrets {
	//	list = append(list, &proto.SecretList{Id: uint32(data.Id), Title: data.Title})
	//}
	//fmt.Println(list)
	//
	//for _, data := range ms.TextSecrets {
	//	list = append(list, &proto.SecretList{Id: uint32(data.Id), Title: data.Title})
	//}
	//fmt.Println(list)
	//
	//for _, data := range ms.CardSecrets {
	//	list = append(list, &proto.SecretList{Id: uint32(data.Id), Title: data.Title})
	//}

	return list
}
