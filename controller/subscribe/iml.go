package subscribe

import (
	"fmt"
	"github.com/gin-gonic/gin"

	"github.com/eolinker/apipark/module/subscribe"
	subscribe_dto "github.com/eolinker/apipark/module/subscribe/dto"
)

var (
	_ ISubscribeController = (*imlSubscribeController)(nil)
)

type imlSubscribeController struct {
	module subscribe.ISubscribeModule `autowired:""`
}

func (i *imlSubscribeController) PartitionServices(ctx *gin.Context, app string) ([]*subscribe_dto.PartitionServiceItem, error) {
	return i.module.PartitionServices(ctx, app)
}

func (i *imlSubscribeController) SearchSubscriptions(ctx *gin.Context, partitionId string, projectId string, keyword string) ([]*subscribe_dto.SubscriptionItem, error) {
	return i.module.SearchSubscriptions(ctx, partitionId, projectId, keyword)
}

func (i *imlSubscribeController) RevokeSubscription(ctx *gin.Context, project string, uuid string) error {
	return i.module.RevokeSubscription(ctx, project, uuid)
}

func (i *imlSubscribeController) DeleteSubscription(ctx *gin.Context, project string, uuid string) error {
	return i.module.DeleteSubscription(ctx, project, uuid)
}

func (i *imlSubscribeController) AddSubscriber(ctx *gin.Context, project string, input *subscribe_dto.AddSubscriber) error {
	return i.module.AddSubscriber(ctx, project, input)
}

func (i *imlSubscribeController) DeleteSubscriber(ctx *gin.Context, project string, serviceId string, applicationId string) error {
	return i.module.DeleteSubscriber(ctx, project, serviceId, applicationId)
}

func (i *imlSubscribeController) RevokeApply(ctx *gin.Context, project string, uuid string) error {
	return i.module.RevokeApply(ctx, project, uuid)
}

func (i *imlSubscribeController) Search(ctx *gin.Context, project string, keyword string) ([]*subscribe_dto.Subscriber, error) {
	return i.module.SearchSubscribers(ctx, project, keyword)
}

var _ ISubscribeApprovalController = (*imlSubscribeApprovalController)(nil)

type imlSubscribeApprovalController struct {
	module subscribe.ISubscribeApprovalModule `autowired:""`
}

func (i *imlSubscribeApprovalController) GetApprovalList(ctx *gin.Context, project string, status int) ([]*subscribe_dto.ApprovalItem, error) {
	return i.module.GetApprovalList(ctx, project, status)
}

func (i *imlSubscribeApprovalController) GetApprovalDetail(ctx *gin.Context, project string, id string) (*subscribe_dto.Approval, error) {
	return i.module.GetApprovalDetail(ctx, project, id)
}

func (i *imlSubscribeApprovalController) Approval(ctx *gin.Context, project string, id string, approveInfo *subscribe_dto.Approve) error {
	switch approveInfo.Operate {
	case "pass":
		return i.module.Pass(ctx, project, id, approveInfo)
	case "refuse":
		return i.module.Reject(ctx, project, id, approveInfo)
	}
	return fmt.Errorf("unknown operate: %s", approveInfo.Operate)
}
