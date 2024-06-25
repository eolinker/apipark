package partition

import (
	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/store"
	"reflect"
)

type IPartitionStore interface {
	store.ISearchStore[Partition]
}

type storePartition struct {
	store.SearchStoreSoftDelete[Partition] // 用struct方式继承,会自动填充并初始化表
}

func init() {
	autowire.Auto[IPartitionStore](func() reflect.Value {
		return reflect.ValueOf(new(storePartition))
	})

}
