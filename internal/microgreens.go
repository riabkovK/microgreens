package internal

type MicrogreensList struct {
	Id                  int    `json:"id"`
	Name                string `json:"name" validate:"required"`
	MicrogreensFamilyId int    `json:"microgreensFamilyId" validate:"required"`
}

type MicrogreensFamily struct {
	Id   int    `json:"id"`
	Name string `json:"name" validate:"required"`
}

type MicrogreensItem struct {
	Id                  int    `json:"id"`
	Name                string `json:"name" validate:"required"`
	MicrogreensFamilyId int    `json:"microgreensFamilyId" validate:"required"`
	Price               int    `json:"price" validate:"required"`
}

type UsersMicrogreensList struct {
	Id                int
	UserId            int
	MicrogreensListId int
}
