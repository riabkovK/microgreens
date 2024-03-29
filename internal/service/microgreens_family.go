package service

import (
	"github.com/riabkovK/microgreens/internal/domain"
	"github.com/riabkovK/microgreens/internal/storage"
)

type MicrogreensFamilyService struct {
	storages storage.MicrogreensFamily
}

func NewMicrogreensFamilyService(storages storage.MicrogreensFamily) *MicrogreensFamilyService {
	return &MicrogreensFamilyService{storages: storages}
}

func (mfs *MicrogreensFamilyService) Create(family domain.MicrogreensFamilyRequest) (int, error) {
	return mfs.storages.Create(family)
}

func (mfs *MicrogreensFamilyService) GetAll() ([]domain.MicrogreensFamily, error) {
	return mfs.storages.GetAll()
}

func (mfs *MicrogreensFamilyService) GetById(familyId int) (domain.MicrogreensFamily, error) {
	return mfs.storages.GetById(familyId)
}

func (mfs *MicrogreensFamilyService) Delete(familyId int) (int, error) {
	return mfs.storages.Delete(familyId)
}

func (mfs *MicrogreensFamilyService) Update(familyId int, request domain.UpdateMicrogreensFamilyRequest) error {
	if err := request.Validate(); err != nil {
		return err
	}
	return mfs.storages.Update(familyId, request)
}
