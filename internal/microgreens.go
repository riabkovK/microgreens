package internal

type MicrogreensList struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	MicrogreensFamilyId int    `json:"family"`
}

type MicrogreensFamily struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type MicrogreensItem struct {
	Id                  int    `json:"id"`
	Name                string `json:"name"`
	MicrogreensFamilyId int    `json:"family"`
	Price               int    `json:"price"`
}

type UsersMicrogreensList struct {
	Id                int
	UserId            int
	MicrogreensListId int
}
