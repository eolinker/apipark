package service

import (
	"context"
	"errors"
	"time"

	"github.com/eolinker/go-common/utils"

	"github.com/eolinker/go-common/auto"

	"gorm.io/gorm"

	"github.com/eolinker/apipark/stores/service"

	"github.com/eolinker/apipark/service/universally"
)

var (
	_ IServiceService    = (*imlServiceService)(nil)
	_ IPartitionsService = (*imlPartitionsService)(nil)
	_ ITagService        = (*imlTagService)(nil)
	_ IDocService        = (*imlDocService)(nil)
	_ IApiService        = (*imlApiService)(nil)
)

type imlServiceService struct {
	store service.IServiceStore `autowired:""`
	universally.IServiceGet[Service]
	universally.IServiceDelete
	universally.IServiceCreate[CreateService]
	universally.IServiceEdit[EditService]
}

func (i *imlServiceService) SearchByUuids(ctx context.Context, keyword string, uuids ...string) ([]*Service, error) {
	w := make(map[string]interface{})
	if len(uuids) > 0 {
		w["uuid"] = uuids
	}
	services, err := i.store.Search(ctx, keyword, w)
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(services, FromEntity), nil
}

func (i *imlServiceService) ListByProject(ctx context.Context, pids ...string) ([]*Service, error) {
	w := make(map[string]interface{})
	if len(pids) > 0 {
		w["project"] = pids
	}
	sorts := []string{"update_at desc", "name asc"}
	services, err := i.store.List(ctx, w, sorts...)
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(services, FromEntity), nil
}

func (i *imlServiceService) SearchServicePage(ctx context.Context, input SearchServicePage, sort ...string) ([]*Service, int64, error) {
	condition := make(map[string]interface{})

	if len(input.Uuids) > 0 {
		condition["uuid"] = input.Uuids
	}
	if len(input.Catalogue) > 0 {
		condition["catalogue"] = input.Catalogue
	}

	if len(sort) == 0 {
		sort = []string{"update_at desc", "name asc"}
	}

	condition["status"] = "on"
	condition["service_type"] = "public"
	items, total, err := i.store.SearchByPage(ctx, input.Keyword, condition, input.Page, input.Size, sort...)
	if err != nil {
		return nil, 0, err
	}
	return utils.SliceToSlice(items, FromEntity), total, nil
}

func (i *imlServiceService) GetLabels(ctx context.Context, ids ...string) map[string]string {
	if len(ids) == 0 {
		return nil
	}
	list, err := i.store.ListQuery(ctx, "`uuid` in (?)", []interface{}{ids}, "id")
	if err != nil {
		return nil
	}
	return utils.SliceToMapO(list, func(i *service.Service) (string, string) {
		return i.UUID, i.Name
	})
}

func (i *imlServiceService) OnComplete() {
	i.IServiceGet = universally.NewGetSoftDelete[Service, service.Service](i.store, FromEntity)
	i.IServiceCreate = universally.NewCreatorSoftDelete[CreateService, service.Service](i.store, "service", createEntityHandler, uniquestHandler, labelHandler)
	i.IServiceDelete = universally.NewSoftDelete[service.Service](i.store)
	i.IServiceEdit = universally.NewEdit[EditService, service.Service](i.store, updateHandler, labelHandler)
	auto.RegisterService("service", i)
}

func labelHandler(e *service.Service) []string {
	return []string{e.Name, e.UUID, e.Tag}
}
func uniquestHandler(i *CreateService) []map[string]interface{} {
	return []map[string]interface{}{{"uuid": i.Uuid}}
}
func createEntityHandler(i *CreateService) *service.Service {
	now := time.Now()
	return &service.Service{
		UUID:        i.Uuid,
		Name:        i.Name,
		Description: i.Description,
		Logo:        i.Logo,
		ServiceType: i.ServiceType,
		Project:     i.Project,
		Team:        i.Team,
		Catalogue:   i.Catalogue,
		Status:      i.Status,
		Tag:         i.Tag,
		CreateAt:    now,
		UpdateAt:    now,
	}
}
func updateHandler(e *service.Service, i *EditService) {
	if i.Name != nil {
		e.Name = *i.Name
	}
	if i.Description != nil {
		e.Description = *i.Description
	}
	if i.Logo != nil {
		e.Logo = *i.Logo
	}
	if i.ServiceType != nil {
		e.ServiceType = *i.ServiceType
	}
	if i.Catalogue != nil {
		e.Catalogue = *i.Catalogue
	}
	if i.Status != nil {
		e.Status = *i.Status
	}
	if i.Tag != nil {
		e.Tag = *i.Tag
	}
	e.UpdateAt = time.Now()
}

type imlPartitionsService struct {
	store service.IServicePartitionStore `autowired:""`
}

func (i *imlPartitionsService) PartitionsByService(ctx context.Context, sids ...string) (map[string][]string, error) {
	condition := make(map[string]interface{})
	if len(sids) > 0 {
		condition["sid"] = sids
	}
	list, err := i.store.List(ctx, condition)
	if err != nil {
		return nil, err
	}
	result := make(map[string][]string)
	for _, v := range list {
		if _, ok := result[v.Sid]; !ok {
			result[v.Sid] = []string{}
		}
		result[v.Sid] = append(result[v.Sid], v.Pid)
	}
	return result, nil
}

func (i *imlPartitionsService) List(ctx context.Context, sid string) ([]*Partition, error) {
	partitions, err := i.store.List(ctx, map[string]interface{}{"sid": sid})
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(partitions, func(p *service.Partition) *Partition {
		return &Partition{
			Pid: p.Pid,
			Sid: p.Sid,
		}
	}), nil

}

func (i *imlPartitionsService) Delete(ctx context.Context, pidList []string, sids []string) error {
	if len(pidList) == 0 && len(sids) == 0 {
		return nil
	}
	conditions := make(map[string]interface{})
	if len(pidList) > 0 {
		conditions["pid"] = pidList
	}
	if len(sids) > 0 {
		conditions["sid"] = sids
	}
	_, err := i.store.DeleteWhere(ctx, conditions)
	return err
}

func (i *imlPartitionsService) Create(ctx context.Context, input *CreatePartition) error {
	return i.store.Insert(ctx, &service.Partition{
		Sid: input.Sid,
		Pid: input.Pid,
	})
}

type imlTagService struct {
	store service.IServiceTagStore `autowired:""`
}

func (i *imlTagService) List(ctx context.Context, sids []string, tids []string) ([]*Tag, error) {
	condition := make(map[string]interface{})
	if len(sids) > 0 {
		condition["sid"] = sids
	}
	if len(tids) > 0 {
		condition["tid"] = tids
	}
	result, err := i.store.List(ctx, condition)
	if err != nil {
		return nil, err
	}
	return utils.SliceToSlice(result, func(s *service.Tag) *Tag {
		return &Tag{
			Tid: s.Tid,
			Sid: s.Sid,
		}
	}), nil
}

func (i *imlTagService) Delete(ctx context.Context, tids []string, sids []string) error {
	if len(tids) == 0 && len(sids) == 0 {
		return nil
	}
	conditions := make(map[string]interface{})
	if len(tids) > 0 {
		conditions["tid"] = tids
	}
	if len(sids) > 0 {
		conditions["sid"] = sids
	}
	_, err := i.store.DeleteWhere(ctx, conditions)
	return err
}

func (i *imlTagService) Create(ctx context.Context, input *CreateTag) error {
	return i.store.Insert(ctx, &service.Tag{
		Sid: input.Sid,
		Tid: input.Tid,
	})
}

type imlDocService struct {
	store service.IServiceDocStore `autowired:""`
}

func (i *imlDocService) Save(ctx context.Context, input *SaveDoc) error {
	info, err := i.Get(ctx, input.Sid)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}
	userID := utils.UserId(ctx)
	if info != nil {
		_, err = i.store.Update(ctx, &service.Doc{
			Id:       info.ID,
			Sid:      input.Sid,
			Doc:      input.Doc,
			CreateAt: info.CreateTime,
			UpdateAt: time.Now(),
			Creator:  info.Creator,
			Updater:  userID,
		})
		return err
	}
	return i.store.Insert(ctx, &service.Doc{
		Sid:      input.Sid,
		Doc:      input.Doc,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
		Creator:  userID,
		Updater:  userID,
	})
}

func (i *imlDocService) Get(ctx context.Context, sid string) (*Doc, error) {
	doc, err := i.store.First(ctx, map[string]interface{}{"sid": sid})
	if err != nil {
		return nil, err
	}
	return &Doc{
		ID:         doc.Id,
		DocID:      doc.Sid,
		Creator:    doc.Creator,
		Updater:    doc.Updater,
		Doc:        doc.Doc,
		UpdateTime: doc.UpdateAt,
		CreateTime: doc.CreateAt,
	}, nil
}

func (i *imlDocService) OnComplete() {
}

type imlApiService struct {
	store service.IServiceApiStore `autowired:""`
}

func (i *imlApiService) CountBySids(ctx context.Context, sids ...string) (map[string]int64, error) {
	wm := make(map[string]interface{})
	if len(sids) > 0 {
		wm["sid"] = sids
	}
	return i.store.CountByGroup(ctx, "", wm, "sid")
}

func (i *imlApiService) Count(ctx context.Context, sid string) (int64, error) {
	return i.store.CountWhere(ctx, map[string]interface{}{"sid": sid})
}

func (i *imlApiService) LastSortIndex(ctx context.Context, sid string) (int, error) {
	result, err := i.store.FirstQuery(ctx, "sid = ?", []interface{}{sid}, "sort desc")
	if err != nil {
		return 0, err
	}
	return result.Sort, nil
}

func (i *imlApiService) Bind(ctx context.Context, sid string, aid string, sort int) error {
	return i.store.Insert(ctx, &service.Api{
		Sid:  sid,
		Aid:  aid,
		Sort: sort,
	})
}

func (i *imlApiService) Unbind(ctx context.Context, sid string, aid string) error {
	_, err := i.store.DeleteWhere(ctx, map[string]interface{}{"sid": sid, "aid": aid})
	return err
}

func (i *imlApiService) List(ctx context.Context, sids ...string) ([]*Api, error) {
	if len(sids) == 0 {
		return nil, nil
	}
	apis, err := i.store.List(ctx, map[string]interface{}{"sid": sids}, "sort asc")
	if err != nil {
		return nil, err
	}
	result := make([]*Api, 0, len(apis))
	for _, api := range apis {
		result = append(result, &Api{
			Aid:  api.Aid,
			Sid:  api.Sid,
			Sort: api.Sort,
		})
	}
	return result, nil
}

func (i *imlApiService) Clear(ctx context.Context, sid string) error {
	_, err := i.store.DeleteWhere(ctx, map[string]interface{}{"sid": sid})
	return err
}
