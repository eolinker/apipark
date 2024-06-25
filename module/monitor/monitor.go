package monitor

import (
	"context"
	"reflect"
	"time"

	"github.com/eolinker/go-common/autowire"

	_ "github.com/eolinker/apipark/module/monitor/driver/influxdb-v2"
	monitor_dto "github.com/eolinker/apipark/module/monitor/dto"
)

type IMonitorStatisticModule interface {
	TopAPIStatistics(ctx context.Context, partitionId string, limit int, input *monitor_dto.CommonInput) ([]*monitor_dto.ApiStatisticItem, error)
	TopProviderStatistics(ctx context.Context, partitionId string, limit int, input *monitor_dto.CommonInput) ([]*monitor_dto.ProjectStatisticItem, error)
	TopSubscriberStatistics(ctx context.Context, partitionId string, limit int, input *monitor_dto.CommonInput) ([]*monitor_dto.ProjectStatisticItem, error)
	// RequestSummary 请求概况
	RequestSummary(ctx context.Context, partitionId string, input *monitor_dto.CommonInput) (*monitor_dto.MonSummaryOutput, error)
	// ProxySummary 转发概况
	ProxySummary(ctx context.Context, partitionId string, input *monitor_dto.CommonInput) (*monitor_dto.MonSummaryOutput, error)

	// InvokeTrend 调用次数趋势
	InvokeTrend(ctx context.Context, partitionId string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error)

	// MessageTrend 消息趋势
	MessageTrend(ctx context.Context, partitionId string, input *monitor_dto.CommonInput) (*monitor_dto.MonMessageTrend, string, error)
	// ProxyTrend 转发趋势

	ApiStatistics(ctx context.Context, partitionId string, input *monitor_dto.StatisticInput) ([]*monitor_dto.ApiStatisticBasicItem, error)

	SubscriberStatistics(ctx context.Context, partitionId string, input *monitor_dto.StatisticInput) ([]*monitor_dto.ProjectStatisticBasicItem, error)

	ProviderStatistics(ctx context.Context, partitionId string, input *monitor_dto.StatisticInput) ([]*monitor_dto.ProjectStatisticBasicItem, error)

	APITrend(ctx context.Context, partitionId string, apiId string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error)

	ProviderTrend(ctx context.Context, partitionId string, providerId string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error)

	SubscriberTrend(ctx context.Context, partitionId string, subscriberId string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error)

	InvokeTrendWithSubscriberAndApi(ctx context.Context, partitionId string, apiId string, subscriberId string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error)
	InvokeTrendWithProviderAndApi(ctx context.Context, partitionId string, providerId string, apiId string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error)

	ProviderStatisticsOnApi(ctx context.Context, partitionId string, apiId string, input *monitor_dto.StatisticInput) ([]*monitor_dto.ProjectStatisticBasicItem, error)
	ApiStatisticsOnProvider(ctx context.Context, partitionId string, providerId string, input *monitor_dto.StatisticInput) ([]*monitor_dto.ApiStatisticBasicItem, error)
	ApiStatisticsOnSubscriber(ctx context.Context, partitionId string, subscriberId string, input *monitor_dto.StatisticInput) ([]*monitor_dto.ApiStatisticBasicItem, error)
	SubscriberStatisticsOnApi(ctx context.Context, partitionId string, apiId string, input *monitor_dto.StatisticInput) ([]*monitor_dto.ProjectStatisticBasicItem, error)
}

func init() {
	autowire.Auto[IMonitorStatisticModule](func() reflect.Value {
		return reflect.ValueOf(new(imlMonitorStatisticModule))
	})
}
func formatTimeByMinute(org int64) time.Time {
	t := time.Unix(org, 0)
	location, _ := time.LoadLocation("Asia/Shanghai")
	return time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), 0, 0, location)
}
