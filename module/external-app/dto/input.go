package external_app_dto

type CreateExternalApp struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type EditExternalApp struct {
	Name *string `json:"name"`
	Desc *string `json:"desc"`
}
