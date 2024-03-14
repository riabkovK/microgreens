package internal

type MicrogreensList struct {
	Id                  int    `json:"id" db:"id"`
	Name                string `json:"name" db:"name" validate:"required"`
	Description         string `json:"description" db:"description"`
	MicrogreensFamilyId int    `json:"microgreensFamilyId" db:"microgreens_family_id"`
}

type MicrogreensFamily struct {
	Id   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name" validate:"required"`
}

type MicrogreensItem struct {
	Id                  int    `json:"id" db:"id" `
	Name                string `json:"name" db:"name" validate:"required"`
	Description         string `json:"description" db:"description"`
	MicrogreensFamilyId int    `json:"microgreensFamilyId" db:"microgreens_family_id"`
	Price               int    `json:"price" db:"price" validate:"required"`
}

type UsersMicrogreensList struct {
	Id                int
	UserId            int
	MicrogreensListId int
}
