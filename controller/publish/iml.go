package publish

import (
	"strconv"

	"github.com/eolinker/apipark/module/publish"
	"github.com/eolinker/apipark/module/publish/dto"
	"github.com/eolinker/apipark/module/release"
	dto2 "github.com/eolinker/apipark/module/release/dto"
	"github.com/gin-gonic/gin"
)

var (
	_ IPublishController = (*imlPublishController)(nil)
)

type imlPublishController struct {
	publishModule publish.IPublishModule `autowired:""`
	releaseModule release.IReleaseModule `autowired:""`
}

func (c *imlPublishController) PublishStatuses(ctx *gin.Context, project string, id string) ([]*dto.PublishStatus, error) {
	return c.publishModule.PublishStatuses(ctx, project, id)
}

func (c *imlPublishController) ApplyOnRelease(ctx *gin.Context, project string, input *dto.ApplyOnReleaseInput) (*dto.Publish, error) {
	newReleaseId, err := c.releaseModule.Create(ctx, project, &dto2.CreateInput{
		Version: input.Version,
		Remark:  input.VersionRemark,
	})
	if err != nil {
		return nil, err
	}
	apply, err := c.publishModule.Apply(ctx, project, &dto.ApplyInput{
		Release: newReleaseId,
		Remark:  input.PublishRemark,
	})
	if err != nil {
		return nil, err
	}
	return apply, nil
}

func (c *imlPublishController) Apply(ctx *gin.Context, project string, input *dto.ApplyInput) (*dto.Publish, error) {
	apply, err := c.publishModule.Apply(ctx, project, input)
	if err != nil {
		return nil, err
	}
	return apply, nil
}

func (c *imlPublishController) CheckPublish(ctx *gin.Context, project string, releaseId string) (*dto.DiffOut, error) {
	return c.publishModule.CheckPublish(ctx, project, releaseId)
}

func (c *imlPublishController) Close(ctx *gin.Context, project string, id string) error {
	err := c.publishModule.Stop(ctx, project, id)
	if err != nil {
		return err
	}
	return nil
}

func (c *imlPublishController) Stop(ctx *gin.Context, project string, id string) error {
	return c.publishModule.Stop(ctx, project, id)
}

func (c *imlPublishController) Refuse(ctx *gin.Context, project string, id string, input *dto.Comments) error {
	return c.publishModule.Refuse(ctx, project, id, input.Comments)
}

func (c *imlPublishController) Accept(ctx *gin.Context, project string, id string, input *dto.Comments) error {
	return c.publishModule.Accept(ctx, project, id, input.Comments)
}

func (c *imlPublishController) Publish(ctx *gin.Context, project string, id string) error {
	return c.publishModule.Publish(ctx, project, id)
}

func (c *imlPublishController) ListPage(ctx *gin.Context, project string, page, pageSize string) ([]*dto.Publish, int, int, int64, error) {
	pageNum, _ := strconv.Atoi(page)
	pageSizeNum, _ := strconv.Atoi(pageSize)
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSizeNum <= 0 {
		pageSizeNum = 50
	}
	list, total, err := c.publishModule.List(ctx, project, pageNum, pageSizeNum)
	if err != nil {
		return nil, 0, 0, 0, err
	}

	return list, pageNum, pageSizeNum, total, nil
}

func (c *imlPublishController) Detail(ctx *gin.Context, project string, id string) (*dto.PublishDetail, error) {
	return c.publishModule.Detail(ctx, project, id)
}
