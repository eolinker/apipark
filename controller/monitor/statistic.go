package monitor

import (
	"reflect"
	"time"

	"github.com/eolinker/go-common/autowire"
	"github.com/gin-gonic/gin"

	monitor_dto "github.com/eolinker/apipark/module/monitor/dto"
)

type IMonitorStatisticController interface {
	Top10(ctx *gin.Context, partition string, input *monitor_dto.Top10Input) (interface{}, error)
	Summary(ctx *gin.Context, partition string, input *monitor_dto.CommonInput) (*monitor_dto.MonSummaryOutput, *monitor_dto.MonSummaryOutput, error)
	OverviewInvokeTrend(ctx *gin.Context, partition string, input *monitor_dto.CommonInput) ([]time.Time, []int64, []int64, []int64, []int64, []float64, []float64, string, error)
	OverviewMessageTrend(ctx *gin.Context, partition string, input *monitor_dto.CommonInput) ([]time.Time, []float64, []float64, string, error)

	Statistics(ctx *gin.Context, partition string, dataType string, input *monitor_dto.StatisticInput) (interface{}, error)

	InvokeTrend(ctx *gin.Context, partition string, dataType string, id string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error)

	InvokeTrendInner(ctx *gin.Context, partition string, dataType string, typ string, api string, provider string, subscriber string, input *monitor_dto.CommonInput) (*monitor_dto.MonInvokeCountTrend, string, error)
	StatisticsInner(ctx *gin.Context, partition string, dataType string, typ string, id string, input *monitor_dto.StatisticInput) (interface{}, error)
}

func init() {
	autowire.Auto[IMonitorStatisticController](func() reflect.Value {
		return reflect.ValueOf(new(imlMonitorStatisticController))
	})
}
