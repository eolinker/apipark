package publish

import (
	"reflect"

	"github.com/eolinker/apipark/module/publish/dto"
	"github.com/eolinker/go-common/autowire"
	"github.com/gin-gonic/gin"
)

var (
	_ IPublishController = (*imlPublishController)(nil)
)

type IPublishController interface {
	CheckPublish(ctx *gin.Context, project string, releaseId string) (*dto.DiffOut, error)

	ApplyOnRelease(ctx *gin.Context, project string, input *dto.ApplyOnReleaseInput) (*dto.Publish, error)
	Apply(ctx *gin.Context, project string, input *dto.ApplyInput) (*dto.Publish, error)
	Close(ctx *gin.Context, project string, id string) error
	Stop(ctx *gin.Context, project string, id string) error
	Refuse(ctx *gin.Context, project string, id string, input *dto.Comments) error
	Accept(ctx *gin.Context, project string, id string, input *dto.Comments) error
	Publish(ctx *gin.Context, project string, id string) error
	ListPage(ctx *gin.Context, project string, page, pageSize string) ([]*dto.Publish, int, int, int64, error)
	Detail(ctx *gin.Context, project string, id string) (*dto.PublishDetail, error)
	PublishStatuses(ctx *gin.Context, project string, id string) ([]*dto.PublishStatus, error)
}

func init() {
	autowire.Auto[IPublishController](func() reflect.Value {
		return reflect.ValueOf(new(imlPublishController))
	})
}
