package organization

import (
	"context"
	"reflect"

	"github.com/eolinker/go-common/auto"
	"github.com/eolinker/go-common/autowire"
)

type IOrganizationService interface {
	Create(ctx context.Context, id, name, description, prefix, master string, partitions []string) (*Organization, error)
	Edit(ctx context.Context, id string, name, description, master *string, partitions *[]string) (*Organization, error)
	Get(ctx context.Context, id string) (*Organization, error)
	Search(ctx context.Context, keyword string) ([]*Organization, error)
	Delete(ctx context.Context, id string) error
	Partitions(ctx context.Context, id string) ([]string, error)
	PartitionsByOrganization(ctx context.Context, orgId ...string) (map[string][]string, error)
	All(ctx context.Context) ([]*Organization, error)
	auto.CompleteService
}

func init() {
	autowire.Auto[IOrganizationService](func() reflect.Value {
		return reflect.ValueOf(new(imlOrganizationService))
	})
}
