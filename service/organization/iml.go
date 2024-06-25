package organization

import (
	"context"
	"errors"
	"time"

	"github.com/eolinker/apipark/stores/organization"
	"github.com/eolinker/go-common/auto"
	"github.com/eolinker/go-common/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	_ IOrganizationService = (*imlOrganizationService)(nil)
)

type imlOrganizationService struct {
	organizationStore organization.IOrganizationStore          `autowired:""`
	partitionStore    organization.IOrganizationPartitionStore `autowired:""`
}

func (s *imlOrganizationService) PartitionsByOrganization(ctx context.Context, orgId ...string) (map[string][]string, error) {
	w := make(map[string]interface{})
	if len(orgId) > 0 {
		w["oid"] = orgId
	}
	partitions, err := s.partitionStore.List(ctx, w)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMapArrayO(partitions, func(p *organization.Partition) (string, string) {
		return p.Oid, p.Pid
	}), nil
}

func (s *imlOrganizationService) All(ctx context.Context) ([]*Organization, error) {
	list, err := s.organizationStore.List(ctx, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(list, func(e *organization.Organization) *Organization {
		return fromEntity(e, nil)
	}), nil
}

func (s *imlOrganizationService) OnComplete() {
	auto.RegisterService("organization", s)
}

func (s *imlOrganizationService) GetLabels(ctx context.Context, ids ...string) map[string]string {
	if len(ids) == 0 {
		return nil
	}
	list, err := s.organizationStore.ListQuery(ctx, "`uuid` in (?)", []interface{}{ids}, "id")
	if err != nil {
		return nil
	}
	return utils.SliceToMapO(list, func(i *organization.Organization) (string, string) {
		return i.UUID, i.Name
	})
}

func (s *imlOrganizationService) Create(ctx context.Context, id, name, description, prefix, master string, partitions []string) (*Organization, error) {
	if id == "" {
		id = uuid.NewString()
	}
	oe, err := s.organizationStore.FirstQuery(ctx, "`uuid`=?", []interface{}{id}, "id")
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if oe != nil {
		return nil, errors.New("organization already exists")
	}
	now := time.Now()
	userId := utils.UserId(ctx)
	ne := &organization.Organization{
		Id:          0,
		UUID:        id,
		Name:        name,
		CreateAt:    now,
		UpdateAt:    now,
		Description: description,
		Master:      master,
		Prefix:      prefix,
		Updater:     userId,
		Creator:     userId,
	}
	pel := make([]*organization.Partition, 0, len(partitions))
	for _, p := range partitions {
		pel = append(pel, &organization.Partition{
			Id:         0,
			Oid:        id,
			Pid:        p,
			CreateTime: ne.CreateAt,
		})
	}
	err = s.partitionStore.Transaction(ctx, func(ctx context.Context) error {
		err := s.organizationStore.Insert(ctx, ne)
		if err != nil {
			return err
		}

		err = s.organizationStore.SetLabels(ctx, ne.Id, ne.Name, ne.Prefix, ne.Master, ne.Description)
		if err != nil {
			return err
		}
		return s.partitionStore.Insert(ctx, pel...)
	})
	if err != nil {
		return nil, err
	}
	return fromEntity(ne, utils.SliceToSlice(pel, func(p *organization.Partition) string { return p.Pid })), nil
}

func (s *imlOrganizationService) Edit(ctx context.Context, id string, name, description, master *string, partitions *[]string) (*Organization, error) {
	if id == "" {
		return nil, errors.New("id is empty")
	}
	operator := utils.UserId(ctx)
	ve, err := s.organizationStore.FirstQuery(ctx, "`uuid`=?", []interface{}{id}, "id")
	if err != nil {
		return nil, err
	}
	if ve == nil {
		return nil, errors.New("organization not exists")
	}
	if name != nil {
		ve.Name = *name
	}
	if description != nil {
		ve.Description = *description
	}
	if master != nil {
		ve.Master = *master
	}
	ve.UpdateAt = time.Now()
	ve.Updater = operator
	pel := make([]*organization.Partition, 0, len(*partitions))
	err = s.organizationStore.Transaction(ctx, func(ctx context.Context) error {
		update, err := s.organizationStore.Update(ctx, ve)
		if err != nil {
			return err
		}
		if update == 0 {
			return errors.New("organization not exists")
		}
		err = s.organizationStore.SetLabels(ctx, ve.Id, ve.Name, ve.Prefix, ve.Master, ve.Description)
		// todo: update labels 需要增加更多参与查询的内容
		if err != nil {
			return err
		}
		if partitions != nil {
			_, err := s.partitionStore.DeleteWhere(ctx, map[string]interface{}{
				"`oid`": id,
			})
			if err != nil {
				return err
			}
			for _, p := range *partitions {
				pel = append(pel, &organization.Partition{
					Id:         0,
					Oid:        id,
					Pid:        p,
					CreateTime: ve.UpdateAt,
				})
			}
			err = s.partitionStore.Insert(ctx, pel...)
		} else {
			pelO, err := s.partitionStore.List(ctx, map[string]interface{}{
				"`oid`": id,
			})
			if err != nil {
				return err
			}
			pel = pelO
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return fromEntity(ve, utils.SliceToSlice(pel, func(p *organization.Partition) string { return p.Pid })), nil
}

func (s *imlOrganizationService) Get(ctx context.Context, id string) (*Organization, error) {
	ev, err := s.organizationStore.FirstQuery(ctx, "`uuid`=?", []interface{}{id}, "id")
	if err != nil {
		return nil, err
	}
	pel, err := s.partitionStore.List(ctx, map[string]interface{}{"`oid`": id})
	if err != nil {
		return nil, err
	}
	return fromEntity(ev, utils.SliceToSlice(pel, func(p *organization.Partition) string { return p.Pid })), nil

}

func (s *imlOrganizationService) Search(ctx context.Context, keyword string) ([]*Organization, error) {
	ev, err := s.organizationStore.Search(ctx, keyword, nil)
	if err != nil {
		return nil, err
	}
	pelAll, err := s.partitionStore.List(ctx, nil)
	if err != nil {
		return nil, err
	}
	pelMap := utils.SliceToMapArray(pelAll, func(i *organization.Partition) string {
		return i.Oid
	})

	return utils.SliceToSlice(ev, func(i *organization.Organization) *Organization {
		return fromEntity(i, utils.SliceToSlice(pelMap[i.UUID], func(p *organization.Partition) string { return p.Pid }))
	}), nil
}

func (s *imlOrganizationService) Delete(ctx context.Context, id string) error {
	_, err := s.organizationStore.FirstQuery(ctx, "`uuid`=?", []interface{}{id}, "id")
	if err != nil {
		return err
	}
	return s.organizationStore.Transaction(ctx, func(ctx context.Context) error {
		deleteCount, err := s.organizationStore.DeleteWhere(ctx, map[string]interface{}{
			"`uuid`": id,
		})
		if err != nil {
			return err
		}
		if deleteCount != 1 {
			return errors.New("delete organization failed")
		}
		_, err = s.partitionStore.DeleteWhere(ctx, map[string]interface{}{
			"`oid`": id,
		})
		return err
	})

}

func (s *imlOrganizationService) Partitions(ctx context.Context, id string) ([]string, error) {
	ps, err := s.partitionStore.List(ctx, map[string]interface{}{"`oid`": id})
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(ps, func(i *organization.Partition) string {
		return i.Pid
	}), nil

}
