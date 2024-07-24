package project_dto

type CreateProject struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	Prefix string `json:"prefix"`
	//Master      string   `json:"master"`
	Partition   []string `json:"partition" aocheck:"partition"`
	Description string   `json:"description"`
	AsApp       *bool    `json:"as_app"`
	AsServer    *bool    `json:"as_server"`
}

type EditProject struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	//Master      *string `json:"master" aocheck:"user"`
}

type CreateApp struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateApp struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type EditMemberRole struct {
	Roles []string `json:"roles"`
}

type Users struct {
	Users []string `json:"users" aocheck:"user"`
}

type EditProjectMember struct {
	Roles []string `json:"roles" aocheck:"role"`
}
