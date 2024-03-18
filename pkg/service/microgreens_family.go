package service

import (
	"github.com/riabkovK/microgreens/internal"
	"github.com/riabkovK/microgreens/internal/storage"
)

type MicrogreensFamilyService struct {
	storages storage.MicrogreensFamily
}

func NewMicrogreensFamilyService(storages storage.MicrogreensFamily) *MicrogreensFamilyService {
	return &MicrogreensFamilyService{storages: storages}
}

func (mfs *MicrogreensFamilyService) Create(family internal.MicrogreensFamily) (int, error) {
	return mfs.storages.Create(family)
}

func (mfs *MicrogreensFamilyService) GetAll() ([]internal.MicrogreensFamily, error) {
	return mfs.storages.GetAll()
}

func (mfs *MicrogreensFamilyService) GetById(familyId int) (internal.MicrogreensFamily, error) {
	return mfs.storages.GetById(familyId)
}

func (mfs *MicrogreensFamilyService) Delete(familyId int) error {
	return mfs.storages.Delete(familyId)
}

func (mfs *MicrogreensFamilyService) Update(familyId int, request internal.UpdateMicrogreensFamilyRequest) error {
	return mfs.storages.Update(familyId, request)
}
