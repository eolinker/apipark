package project

import (
	"time"

	"github.com/eolinker/apipark/stores/project"
)

type Project struct {
	Id          string
	Name        string
	Description string
	Team        string
	//Master      string
	Prefix     string
	AsServer   bool
	AsApp      bool
	CreateTime time.Time
	UpdateTime time.Time
}

func FromEntity(e *project.Project) *Project {
	return &Project{
		Id:          e.UUID,
		Name:        e.Name,
		Description: e.Description,
		Team:        e.Team,
		//Master:      e.Master,
		Prefix:     e.Prefix,
		AsServer:   e.AsServer,
		AsApp:      e.AsApp,
		CreateTime: e.CreateAt,
		UpdateTime: e.UpdateAt,
	}
}

type CreateProject struct {
	Id          string
	Name        string
	Description string
	//Master      string
	Team     string
	Prefix   string
	AsServer bool
	AsApp    bool
}

type EditProject struct {
	Name        *string
	Description *string
	//Master      *string
}

type Partition struct {
	Partition string
	Project   string
}
