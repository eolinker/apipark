package partition

import (
	"context"
	"time"

	"github.com/eolinker/go-common/utils"

	"github.com/eolinker/apipark/service/universally"
	"github.com/eolinker/apipark/stores/partition"
	"github.com/eolinker/go-common/auto"
)

var (
	_ IPartitionService = (*imlPartitionService)(nil)
)

type imlPartitionService struct {
	store partition.IPartitionStore `autowired:""`
	universally.IServiceGet[Partition]
	universally.IServiceDelete
	universally.IServiceCreate[CreatePartition]
	universally.IServiceEdit[EditPartition]
}

func (p *imlPartitionService) GetLabels(ctx context.Context, ids ...string) map[string]string {
	if len(ids) == 0 {
		return nil
	}
	list, err := p.store.ListQuery(ctx, "`uuid` in (?)", []interface{}{ids}, "id")
	if err != nil {
		return nil
	}
	return utils.SliceToMapO(list, func(i *partition.Partition) (string, string) {
		return i.UUID, i.Name
	})
}

func (p *imlPartitionService) OnComplete() {
	p.IServiceGet = universally.NewGetSoftDelete[Partition, partition.Partition](p.store, toModel)

	p.IServiceDelete = universally.NewSoftDelete[partition.Partition](p.store)

	p.IServiceCreate = universally.NewCreatorSoftDelete[CreatePartition, partition.Partition](p.store, "partition", createEntityHandler, uniquestHandler, labelHandler)

	p.IServiceEdit = universally.NewEdit[EditPartition, partition.Partition](p.store, updateHandler, labelHandler)

	auto.RegisterService("partition", p)
}
func uniquestHandler(i *CreatePartition) []map[string]interface{} {
	return []map[string]interface{}{{"uuid": i.Uuid}, {"name": i.Name}}

}
func updateHandler(e *partition.Partition, i *EditPartition) {
	if i.Name != nil {
		e.Name = *i.Name
	}
	if i.Resume != nil {
		e.Resume = *i.Resume
	}
	if i.Url != nil {
		e.Url = *i.Url
	}
	if i.Prefix != nil {
		e.Prefix = *i.Prefix
	}
	e.UpdateAt = time.Now()
}
func createEntityHandler(i *CreatePartition) *partition.Partition {
	return &partition.Partition{
		Id:       0,
		UUID:     i.Uuid,
		Name:     i.Name,
		CreateAt: time.Now(),
		UpdateAt: time.Now(),
		Resume:   i.Resume,
		Prefix:   i.Prefix,
		Url:      i.Url,
		Cluster:  i.Cluster,
	}
}

func labelHandler(e *partition.Partition) []string {
	return []string{e.UUID, e.Name, e.Url, e.Prefix, e.Resume}
}

func toModel(i *partition.Partition) *Partition {
	return &Partition{
		UUID:       i.UUID,
		Name:       i.Name,
		Resume:     i.Resume,
		Prefix:     i.Prefix,
		Url:        i.Url,
		Updater:    i.Updater,
		UpdateTime: i.UpdateAt,
		Creator:    i.Creator,
		CreateTime: i.CreateAt,
		Cluster:    i.Cluster,
	}
}
