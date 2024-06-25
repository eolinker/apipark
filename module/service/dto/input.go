package service_dto

type CreateService struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Logo        string   `json:"logo"`
	Tags        []string `json:"tags"`
	ServiceType string   `json:"service_type"`
	Partition   []string `json:"partition" aocheck:"partition"`
	Catalogue   *string  `json:"group" aocheck:"catalogue"`
}

type EditService struct {
	Name        *string  `json:"name"`
	Description *string  `json:"description"`
	Logo        *string  `json:"logo"`
	Tags        []string `json:"tags"`
	ServiceType *string  `json:"service_type"`
	Catalogue   *string  `json:"group" aocheck:"catalogue"`
}

type BindApis struct {
	Apis []string `json:"apis" aocheck:"api"`
}
