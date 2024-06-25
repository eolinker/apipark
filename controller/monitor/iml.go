package monitor

import (
	"fmt"
	"time"

	"github.com/eolinker/apipark/module/monitor"
	monitor_dto "github.com/eolinker/apipark/module/monitor/dto"
	"github.com/gin-gonic/gin"
)

var (
	_ IMonitorStatisticController = (*imlMonitorStatisticController)(nil)
)

type imlMonitorStatisticController struct {
	module monitor.IMonitorStatisticModule `autowired:""`
}

func (i *imlMonitorStatisticController) InvokeTrendInner(ctx *gin.Context, partition string, dataType string, typ string, api string, provider string, subscriber string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error) {
	if dataType == monitor_dto.DataTypeApi && typ == monitor_dto.DataTypeSubscriber || dataType == monitor_dto.DataTypeSubscriber && typ == monitor_dto.DataTypeApi {
		return i.module.InvokeTrendWithSubscriberAndApi(ctx, partition, api, subscriber, input)
	} else if dataType == monitor_dto.DataTypeApi && typ == monitor_dto.DataTypeProvider || dataType == monitor_dto.DataTypeProvider && typ == monitor_dto.DataTypeApi {
		return i.module.InvokeTrendWithProviderAndApi(ctx, partition, provider, api, input)
	}
	return nil, "", fmt.Errorf("unsupported detail type: %s, data type is %s", typ, dataType)
}

func (i *imlMonitorStatisticController) StatisticsInner(ctx *gin.Context, partition string, dataType string, typ string, id string, input *monitor_dto.StatisticInput) (interface{}, error) {
	switch dataType {
	case monitor_dto.DataTypeApi:
		switch typ {
		case monitor_dto.DataTypeProvider:
			return i.module.ProviderStatisticsOnApi(ctx, partition, id, input)
		case monitor_dto.DataTypeSubscriber:
			return i.module.SubscriberStatisticsOnApi(ctx, partition, id, input)
		default:
			return nil, fmt.Errorf("unsupported detail type: %s, data type is %s", typ, dataType)
		}
	case monitor_dto.DataTypeProvider:
		switch typ {
		case monitor_dto.DataTypeApi:
			return i.module.ApiStatisticsOnProvider(ctx, partition, id, input)
		default:
			return nil, fmt.Errorf("unsupported detail type: %s, data type is %s", typ, dataType)
		}
	case monitor_dto.DataTypeSubscriber:
		switch typ {
		case monitor_dto.DataTypeApi:
			return i.module.ApiStatisticsOnSubscriber(ctx, partition, id, input)
		default:
			return nil, fmt.Errorf("unsupported detail type: %s, data type is %s", typ, dataType)
		}
	}
	return nil, fmt.Errorf("unsupported data type: %s", dataType)

}

func (i *imlMonitorStatisticController) InvokeTrend(ctx *gin.Context, partition string, dataType string, id string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error) {
	switch dataType {
	case monitor_dto.DataTypeApi:
		return i.module.APITrend(ctx, partition, id, input)
	case monitor_dto.DataTypeProvider:
		return i.module.ProviderTrend(ctx, partition, id, input)
	case monitor_dto.DataTypeSubscriber:
		return i.module.SubscriberTrend(ctx, partition, id, input)
	default:
		return nil, "", fmt.Errorf("unsupported data type: %s", dataType)
	}
}

func (i *imlMonitorStatisticController) Statistics(ctx *gin.Context, partition string, dataType string, input *monitor_dto.StatisticInput) (interface{}, error) {
	switch dataType {
	case monitor_dto.DataTypeApi:
		return i.module.ApiStatistics(ctx, partition, input)
	case monitor_dto.DataTypeProvider:
		return i.module.ProviderStatistics(ctx, partition, input)
	case monitor_dto.DataTypeSubscriber:
		return i.module.SubscriberStatistics(ctx, partition, input)
	default:
		return nil, fmt.Errorf("unsupported data type: %s", dataType)
	}
}

func (i *imlMonitorStatisticController) OverviewMessageTrend(ctx *gin.Context, partition string, input *monitor_dto.CommonInput) ([]time.Time, []float64, []float64, string, error) {
	trend, timeInterval, err := i.module.MessageTrend(ctx, partition, input)
	if err != nil {
		return nil, nil, nil, "", err
	}

	return trend.Dates, trend.ReqMessage, trend.RespMessage, timeInterval, nil
}

func (i *imlMonitorStatisticController) OverviewInvokeTrend(ctx *gin.Context, partition string, input *monitor_dto.CommonInput) ([]time.Time, []int64, []int64, []int64, []int64, []float64, []float64, string, error) {
	trend, timeInterval, err := i.module.InvokeTrend(ctx, partition, input)
	if err != nil {
		return nil, nil, nil, nil, nil, nil, nil, "", err
	}

	return trend.Date, trend.RequestTotal, trend.ProxyTotal, trend.Status4XX, trend.Status5XX, trend.RequestRate, trend.ProxyRate, timeInterval, nil
}

func (i *imlMonitorStatisticController) Summary(ctx *gin.Context, partition string, input *monitor_dto.CommonInput) (*monitor_dto.MonSummaryOutput, *monitor_dto.MonSummaryOutput, error) {
	requestSummary, err := i.module.RequestSummary(ctx, partition, input)
	if err != nil {
		return nil, nil, err
	}
	proxySummary, err := i.module.ProxySummary(ctx, partition, input)
	if err != nil {
		return nil, nil, err
	}
	return requestSummary, proxySummary, nil
}

func (i *imlMonitorStatisticController) Top10(ctx *gin.Context, partition string, input *monitor_dto.Top10Input) (interface{}, error) {
	switch input.DataType {
	case monitor_dto.DataTypeApi:
		return i.module.TopAPIStatistics(ctx, partition, 10, input.CommonInput)
	case monitor_dto.DataTypeProvider:
		return i.module.TopProviderStatistics(ctx, partition, 10, input.CommonInput)
	case monitor_dto.DataTypeSubscriber:
		return i.module.TopSubscriberStatistics(ctx, partition, 10, input.CommonInput)
	default:
		return nil, fmt.Errorf("unsupported data type: %s", input.DataType)
	}
}
