package organization

import (
	"context"
	"strings"

	"github.com/eolinker/apipark/service/team"

	organization_dto "github.com/eolinker/apipark/module/organization/dto"
	"github.com/eolinker/apipark/service/organization"
	"github.com/eolinker/apipark/service/partition"
	"github.com/eolinker/go-common/auto"
	"github.com/eolinker/go-common/utils"
	"github.com/google/uuid"
)

var (
	_ IOrganizationModule = (*implOrganizationModule)(nil)
)

type implOrganizationModule struct {
	organizationService organization.IOrganizationService `autowired:""`
	partitionService    partition.IPartitionService       `autowired:""`
	teamService         team.ITeamService                 `autowired:""`
}

func (m *implOrganizationModule) Simple(ctx context.Context) ([]*organization_dto.Simple, error) {
	list, err := m.organizationService.All(ctx)
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(list, func(i *organization.Organization) *organization_dto.Simple {
		return &organization_dto.Simple{
			Id:   i.UUID,
			Name: i.Name,
		}
	}), nil
}

func (m *implOrganizationModule) Create(ctx context.Context, input *organization_dto.CreateOrganization) (*organization_dto.Detail, error) {
	id := input.Id
	if id == "" {
		id = uuid.NewString()
	}
	input.Prefix = strings.Trim(strings.Trim(input.Prefix, " "), "/")
	o, err := m.organizationService.Create(ctx, id, input.Name, input.Description, input.Prefix, input.Master, input.Partitions)
	if err != nil {
		return nil, err
	}
	detail := organization_dto.NewDetail(o)
	return detail, nil
}

func (m *implOrganizationModule) Edit(ctx context.Context, id string, input *organization_dto.EditOrganization) (*organization_dto.Detail, error) {
	nv, err := m.organizationService.Edit(ctx, id, input.Name, input.Description, input.Master, input.Partitions)
	if err != nil {
		return nil, err
	}
	detail := organization_dto.NewDetail(nv)

	return detail, nil
}

func (m *implOrganizationModule) Get(ctx context.Context, id string) (*organization_dto.Detail, error) {
	v, err := m.organizationService.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	detail := organization_dto.NewDetail(v)

	return detail, nil
}

func (m *implOrganizationModule) Search(ctx context.Context, keyword string) ([]*organization_dto.Item, error) {
	list, err := m.organizationService.Search(ctx, keyword)
	if err != nil {
		return nil, err
	}
	countMap, err := m.teamService.CountByGroup(ctx, "", nil, "organization")
	if err != nil {
		return nil, err
	}
	itemList := utils.SliceToSlice(list, func(o *organization.Organization) *organization_dto.Item {
		return &organization_dto.Item{
			Id:          o.UUID,
			Name:        o.Name,
			Description: o.Description,
			Master:      auto.UUID(o.Master),
			Partition:   auto.List(o.Partitions),
			Prefix:      o.Prefix,
			CreateTime:  auto.TimeLabel(o.CreateTime),
			UpdateTime:  auto.TimeLabel(o.UpdateTime),
			Updater:     auto.UUID(o.Updater),
			Creator:     auto.UUID(o.Creator),
			CanDelete:   countMap[o.UUID] == 0,
		}
	})
	return itemList, nil
}

func (m *implOrganizationModule) Delete(ctx context.Context, id string) (string, error) {
	err := m.organizationService.Delete(ctx, id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (m *implOrganizationModule) Partitions(ctx context.Context, id string) ([]*organization_dto.Partition, error) {
	partitions, err := m.organizationService.Partitions(ctx, id)
	if err != nil {
		return nil, err
	}
	list, err := m.partitionService.List(ctx, partitions...)
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(list, func(p *partition.Partition) *organization_dto.Partition {
		return &organization_dto.Partition{
			Id:   p.UUID,
			Name: p.Name,
		}
	}), nil

}
