package service

import (
	"github.com/riabkovK/microgreens/internal"
	"github.com/riabkovK/microgreens/internal/storage"
)

type MicrogreensItemService struct {
	storages    storage.MicrogreensItem
	listStorage storage.MicrogreensList
}

func NewMicrogreensItemService(storages storage.MicrogreensItem, listStorage storage.MicrogreensList) *MicrogreensItemService {
	return &MicrogreensItemService{storages: storages, listStorage: listStorage}
}

func (mis *MicrogreensItemService) Create(userId, listId int, microgreensItem internal.MicrogreensItem) (int, error) {
	_, err := mis.listStorage.GetById(userId, listId)
	if err != nil {
		// list does not exist or does not belong to user
		return 0, err
	}
	return mis.storages.Create(listId, microgreensItem)
}

func (mis *MicrogreensItemService) GetAll(userId, listId int) ([]internal.MicrogreensItem, error) {
	return mis.storages.GetAll(userId, listId)
}

func (mis *MicrogreensItemService) GetById(userId, itemId int) (internal.MicrogreensItem, error) {
	return mis.storages.GetById(userId, itemId)
}

func (mis *MicrogreensItemService) Delete(userId, listId int) error {
	return mis.storages.Delete(userId, listId)
}

func (mis *MicrogreensItemService) Update(userId, itemId int, request internal.UpdateMicrogreensItemRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	return mis.storages.Update(userId, itemId, request)
}
