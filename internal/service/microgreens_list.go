package service

import (
	"github.com/riabkovK/microgreens/internal/domain"
	"github.com/riabkovK/microgreens/internal/storage"
)

type MicrogreensListService struct {
	storages storage.MicrogreensList
}

func NewMicrogreensListService(storages storage.MicrogreensList) *MicrogreensListService {
	return &MicrogreensListService{storages: storages}
}

func (mls *MicrogreensListService) Create(userId int, list domain.MicrogreensListRequest) (int, error) {
	return mls.storages.Create(userId, list)
}

func (mls *MicrogreensListService) GetAll(userId int) ([]domain.MicrogreensList, error) {
	return mls.storages.GetAll(userId)
}

func (mls *MicrogreensListService) GetById(userId, listId int) (domain.MicrogreensList, error) {
	return mls.storages.GetById(userId, listId)
}

func (mls *MicrogreensListService) Delete(userId, listId int) (int, error) {
	return mls.storages.Delete(userId, listId)
}

func (mls *MicrogreensListService) Update(userId, listId int, request domain.UpdateMicrogreensListRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	return mls.storages.Update(userId, listId, request)
}
