package upstream

import (
	"context"
	"errors"
	"fmt"

	"github.com/eolinker/apipark/service/partition"

	"github.com/eolinker/go-common/utils"

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
	projectService   project.IProjectService     `autowired:""`
	partitionService partition.IPartitionService `autowired:""`
	upstreamService  upstream.IUpstreamService   `autowired:""`
	transaction      store.ITransaction          `autowired:""`
}

func (i *imlUpstreamModule) Get(ctx context.Context, pid string) (upstream_dto.UpstreamConfig, error) {
	_, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)
	if err != nil {
		return nil, err
	}
	item, err := i.upstreamService.Get(ctx, pid)
	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		return nil, nil
	}
	result := make(map[string]*upstream_dto.Upstream)
	for _, partitionId := range item.Partitions {
		commit, err := i.upstreamService.LatestCommit(ctx, pid, partitionId)
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, err
			}
			continue
		}
		result[partitionId] = upstream_dto.FromClusterConfig(commit.Data)
	}

	return result, nil
}

func (i *imlUpstreamModule) Save(ctx context.Context, pid string, upstreamConfig upstream_dto.UpstreamConfig) (upstream_dto.UpstreamConfig, error) {
	pInfo, err := i.projectService.CheckProject(ctx, pid, projectRuleMustServer)
	if err != nil {
		return nil, err
	}
	partitions, err := i.partitionService.List(ctx)
	if err != nil {
		return nil, err
	}

	partitionMap := utils.SliceToMapO(partitions, func(p *partition.Partition) (string, struct{}) {
		return p.UUID, struct{}{}
	})
	err = i.transaction.Transaction(ctx, func(ctx context.Context) error {
		partitionIds := make([]string, 0, len(upstreamConfig))
		for id, cfg := range upstreamConfig {
			if _, ok := partitionMap[id]; !ok {
				continue
			}
			err = i.upstreamService.SaveCommit(ctx, pid, id, upstream_dto.ConvertUpstream(cfg))
			if err != nil {
				return err
			}

			partitionIds = append(partitionIds, id)
		}
		return i.upstreamService.Save(ctx, &upstream.SaveUpstream{
			UUID:       pid,
			Name:       fmt.Sprintf("upstream-%s", pid),
			Project:    pid,
			Team:       pInfo.Team,
			Partitions: partitionIds,
		})

	})
	if err != nil {
		return nil, err
	}
	return i.Get(ctx, pid)
}
