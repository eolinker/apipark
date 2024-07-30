package upstream

import (
	"context"
	"errors"
	"fmt"

	"github.com/eolinker/apipark/service/cluster"

	"gorm.io/gorm"

	"github.com/eolinker/apipark/service/project"

	"github.com/eolinker/apipark/service/upstream"

	"github.com/eolinker/go-common/store"

	upstream_dto "github.com/eolinker/apipark/module/upstream/dto"
)

var (
	_                     IUpstreamModule = (*imlUpstreamModule)(nil)
	projectRuleMustServer                 = map[string]bool{
		"as_server": true,
	}
)

type imlUpstreamModule struct {
	projectService  project.IProjectService   `autowired:""`
	upstreamService upstream.IUpstreamService `autowired:""`
	transaction     store.ITransaction        `autowired:""`
}

func (i *imlUpstreamModule) Get(ctx context.Context, pid string) (upstream_dto.UpstreamConfig, error) {
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)
	if err != nil {
		return nil, err
	}
	_, err = i.upstreamService.Get(ctx, pid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, nil
	}
	commit, err := i.upstreamService.LatestCommit(ctx, pid, cluster.DefaultClusterID)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, nil
	}

	return upstream_dto.FromClusterConfig(commit.Data), nil
}

func (i *imlUpstreamModule) Save(ctx context.Context, pid string, upstreamConfig upstream_dto.UpstreamConfig) (upstream_dto.UpstreamConfig, error) {
	pInfo, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)
	if err != nil {
		return nil, err
	}

	err = i.transaction.Transaction(ctx, func(ctx context.Context) error {
		err = i.upstreamService.SaveCommit(ctx, pid, cluster.DefaultClusterID, upstream_dto.ConvertUpstream(upstreamConfig))
		if err != nil {
			return err
		}

		return i.upstreamService.Save(ctx, &upstream.SaveUpstream{
			UUID:    pid,
			Name:    fmt.Sprintf("upstream-%s", pid),
			Project: pid,
			Team:    pInfo.Team,
		})

	})
	if err != nil {
		return nil, err
	}
	return i.Get(ctx, pid)
}
