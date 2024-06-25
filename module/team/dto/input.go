package team_dto

type CreateTeam struct {
	Id             string `json:"id"`
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	OrganizationId string `json:"organization" aocheck:"organization"`
	Master         string `json:"master" aocheck:"user"`
}
type EditTeam struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Master      *string `json:"master" aocheck:"user"`
}
