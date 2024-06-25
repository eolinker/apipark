package gateway

import (
	"context"
)

var (
	initHandlers []InitHandler
)

func RegisterInitHandler(handle InitHandler) {
	initHandlers = append(initHandlers, handle)
}
func RegisterInitHandleFunc(handleFunc InitHandleFunc) {
	initHandlers = append(initHandlers, handleFunc)
}

type InitHandleFunc func(ctx context.Context, partitionId string, client IClientDriver) error

func (f InitHandleFunc) Init(ctx context.Context, partitionId string, client IClientDriver) error {
	return f(ctx, partitionId, client)
}

type InitHandler interface {
	Init(ctx context.Context, partitionId string, client IClientDriver) error
}

func InitGateway(ctx context.Context, partitionId string, client IClientDriver) (err error) {
	//defer func() {
	//	if err == nil {
	//		err = client.Commit(ctx)
	//	} else {
	//		errRollback := client.Rollback(ctx)
	//		if errRollback != nil {
	//			log.Warn(err)
	//		}
	//	}
	//}()
	for _, h := range initHandlers {
		err = h.Init(ctx, partitionId, client)
		if err != nil {
			return
		}
	}
	return
}
