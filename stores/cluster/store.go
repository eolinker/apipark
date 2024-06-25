package cluster

import (
	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/store"
	"reflect"
)

type IClusterNodeStore interface {
	store.IBaseStore[ClusterNode]
}
type storeClusterNode struct {
	store.Store[ClusterNode] // 用struct方式继承,会自动填充并初始化表
}
type IClusterNodeAddressStore interface {
	store.IBaseStore[ClusterNodeAddr]
}
type storeClusterNodeAddr struct {
	store.Store[ClusterNodeAddr] // 用struct方式继承,会自动填充并初始化表
}

func init() {
	autowire.Auto[IClusterStore](func() reflect.Value {
		return reflect.ValueOf(new(storeCluster))
	})

	autowire.Auto[IClusterNodeStore](func() reflect.Value {
		return reflect.ValueOf(new(storeClusterNode))
	})

	autowire.Auto[IClusterNodeAddressStore](func() reflect.Value {
		return reflect.ValueOf(new(storeClusterNodeAddr))
	})
}
