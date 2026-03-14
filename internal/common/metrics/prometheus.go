package metrics

import (
	"strings"

	"github.com/phrara/mallive/common/decorator"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	// "github.com/prometheus/client_golang/prometheus/promauto"
)

var _ decorator.MetricsClient = (*PrometheusMetricsClient)(nil)



var (
	// 1. 定义一个计数器（Counter）：记录服务成功请求次数
	successRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "success_requests_total",
			Help: "Total number of successful requests to the order service.",
		},
		[]string{"CQRS", "Action"}, // 定义标签（Labels）
	)
	
	// 2. 定义一个计数器（Counter）：记录服务失败请求次数
	failRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "failed_requests_total",
			Help: "Total number of failed requests to the order service.",
		},
		[]string{"CQRS", "Action"}, // 定义标签（Labels）
	)

	// 3. 定义一个直方图（Histogram）：记录处理耗时
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "inventory_request_duration_seconds",
			Help:    "Histogram of request durations.",
			Buckets: prometheus.DefBuckets, // 默认的耗时分布区间
		},
		[]string{"CQRS", "Action"},
	)

	// "github.com/prometheus/client_golang/prometheus/promauto"
	// 如果你在定义指标时使用了 promauto 包，
	// 如 httpRequests = promauto.NewCounterVec(...)
	// 它会自动将指标注册到全局默认注册表 prometheus.DefaultRegisterer 中。

)

type PrometheusMetricsClient struct {
	registry *prometheus.Registry
}

func NewPrometheusMetricsClient(serviceName string) *PrometheusMetricsClient {
	p := &PrometheusMetricsClient{}
	p.initPrometheus(serviceName)

	return p
}

func (p *PrometheusMetricsClient) GetPromRegistry() *prometheus.Registry {
	return p.registry
}

func (p *PrometheusMetricsClient) initPrometheus(serviceName string)  {
	p.registry = prometheus.NewRegistry()
    
    // 关键：创建一个带全局标签的“包装注册器”
    wrappedReg := prometheus.WrapRegistererWith(prometheus.Labels{
        "service": serviceName, // 建议统一叫 service 而不是 serviceName
    }, p.registry)

    // 使用包装后的注册器进行 MustRegister
    wrappedReg.MustRegister(
        collectors.NewGoCollector(),
        collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
        successRequestsTotal,
        failRequestsTotal,
        requestDuration,
    )

}

func (*PrometheusMetricsClient) Inc(key string, value any) {
	cqrs := strings.Split(key, ".")[0]
	action := strings.Split(key, ".")[1]
	metric := strings.Split(key, ".")[2]
	switch metric {
		// 更新指标
	case "duration":
		requestDuration.WithLabelValues(cqrs, action).Observe(value.(float64))
	case "success":
		successRequestsTotal.WithLabelValues(cqrs, action).Inc()
	case "failure":
		failRequestsTotal.WithLabelValues(cqrs, action).Inc()
	default:
		
	}
}