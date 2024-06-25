package core

import (
	"net/http"

	"github.com/eolinker/go-common/pm3"
)

func (p *plugin) monitorStatisticApis() []pm3.Api {
	return []pm3.Api{
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/monitor/overview/top10", []string{"context", "query:partition", "body"}, []string{"top10"}, p.monitorStatisticController.Top10),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/monitor/overview/summary", []string{"context", "query:partition", "body"}, []string{"request_summary", "proxy_summary"}, p.monitorStatisticController.Summary),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/monitor/overview/invoke", []string{"context", "query:partition", "body"}, []string{"date", "request_total", "proxy_total", "status_4xx", "status_5xx", "request_rate", "proxy_rate", "time_interval"}, p.monitorStatisticController.OverviewInvokeTrend),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/monitor/overview/message", []string{"context", "query:partition", "body"}, []string{"date", "request_message", "response_message", "time_interval"}, p.monitorStatisticController.OverviewMessageTrend),

		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/monitor/:data_type", []string{"context", "query:partition", "rest:data_type", "body"}, []string{"statistics"}, p.monitorStatisticController.Statistics),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/monitor/:data_type/trend", []string{"context", "query:partition", "rest:data_type", "query:id", "body"}, []string{"tendency", "time_interval"}, p.monitorStatisticController.InvokeTrend),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/monitor/:data_type/trend/:typ", []string{"context", "query:partition", "rest:data_type", "rest:typ", "query:api", "query:provider", "query:subscriber", "body"}, []string{"tendency", "time_interval"}, p.monitorStatisticController.InvokeTrendInner),
		pm3.CreateApiWidthDoc(http.MethodPost, "/api/v1/monitor/:data_type/statistics/:typ", []string{"context", "query:partition", "rest:data_type", "rest:typ", "query:id", "body"}, []string{"statistics"}, p.monitorStatisticController.StatisticsInner),
	}
}
