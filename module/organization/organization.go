package organization

import (
	"context"
	organization_dto "github.com/eolinker/apipark/module/organization/dto"
	"github.com/eolinker/go-common/autowire"
	"reflect"
)

type IOrganizationModule interface {
	Create(ctx context.Context, input *organization_dto.CreateOrganization) (*organization_dto.Detail, error)
	Edit(ctx context.Context, id string, input *organization_dto.EditOrganization) (*organization_dto.Detail, error)
	Get(ctx context.Context, id string) (*organization_dto.Detail, error)
	Search(ctx context.Context, keyword string) ([]*organization_dto.Item, error)
	Delete(ctx context.Context, id string) (string, error)
	Simple(ctx context.Context) ([]*organization_dto.Simple, error)
	Partitions(ctx context.Context, id string) ([]*organization_dto.Partition, error)
}

func init() {
	autowire.Auto[IOrganizationModule](func() reflect.Value {
		return reflect.ValueOf(new(implOrganizationModule))
	})
}
