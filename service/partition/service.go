package partition

import (
	"github.com/eolinker/apipark/service/universally"
	"github.com/eolinker/go-common/autowire"
	"reflect"
)

func init() {
	autowire.Auto[IPartitionService](func() reflect.Value {
		return reflect.ValueOf(new(imlPartitionService))
	})
}

type IPartitionService interface {
	universally.IServiceGet[Partition]
	universally.IServiceDelete
	universally.IServiceCreate[CreatePartition]
	universally.IServiceEdit[EditPartition]
}
