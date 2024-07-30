package service

import (
	"context"
	"errors"

	"github.com/eolinker/apipark/service/project"
	"gorm.io/gorm"

	"github.com/eolinker/go-common/auto"

	serviceDto "github.com/eolinker/apipark/module/service/dto"
	"github.com/eolinker/apipark/service/service"
	"github.com/eolinker/go-common/store"
)

var (
	_                     IServiceModule = (*imlServiceModule)(nil)
	projectRuleMustServer                = map[string]bool{
		"as_server": true,
	}
)

type imlServiceModule struct {
	projectService    project.IProjectService `autowired:""`
	serviceDocService service.IDocService     `autowired:""`

	transaction store.ITransaction `autowired:""`
}

func (i *imlServiceModule) ServiceDoc(ctx context.Context, pid string) (*serviceDto.ServiceDoc, error) {
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)

	if err != nil {
		return nil, err
	}
	info, err := i.projectService.Get(ctx, pid)
	if err != nil {
		return nil, err
	}
	doc, err := i.serviceDocService.Get(ctx, pid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return &serviceDto.ServiceDoc{
			Id:   pid,
			Name: info.Name,
			Doc:  "",
		}, nil
	}
	return &serviceDto.ServiceDoc{
		Id:         pid,
		Name:       info.Name,
		Doc:        doc.Doc,
		Creator:    auto.UUID(doc.Creator),
		CreateTime: auto.TimeLabel(doc.CreateTime),
		Updater:    auto.UUID(doc.Updater),
		UpdateTime: auto.TimeLabel(doc.UpdateTime),
	}, nil
}

func (i *imlServiceModule) SaveServiceDoc(ctx context.Context, pid string, input *serviceDto.SaveServiceDoc) error {
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)

	if err != nil {
		return err
	}
	return i.serviceDocService.Save(ctx, &service.SaveDoc{
		Sid: pid,
		Doc: input.Doc,
	})
}
