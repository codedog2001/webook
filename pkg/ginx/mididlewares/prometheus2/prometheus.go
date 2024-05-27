package prometheus2

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

// count只能统计一些只增不减的数，比如http请求的总数，错误数，成功数等
// gauge可以统计一些增加或者减少的量，比如内存，cpu等
//Summary主要用于描述指标的分布情况，例如请求响应时间。 也就是多维度的
//通过SummaryVec，可以在同一个指标名称下使用不同的标签值创建多个度量，从而更全面地了解指标的分布情况。

type Builder struct {
	Namespace  string
	Subsystem  string
	Name       string
	InstanceId string
	Help       string
}

func (b *Builder) BuildResponseTime() gin.HandlerFunc {
	//pattern是指命中的路由
	labels := []string{"method", "path", "status"}
	vector := prometheus.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Help:      b.Help,
		Name:      b.Name + "_rest_time",
		ConstLabels: map[string]string{
			"instance_id": b.InstanceId,
		},
		Objectives: map[float64]float64{
			0.5:   0.01,
			0.9:   0.01,
			0.99:  0.001,
			0.999: 0.0001,
		},
	}, labels)
	prometheus.MustRegister(vector)
	return func(ctx *gin.Context) {
		start := time.Now()
		defer func() {
			duration := time.Since(start).Milliseconds()
			method := ctx.Request.Method
			pattern := ctx.FullPath()
			status := ctx.Writer.Status()
			vector.WithLabelValues(method, pattern, strconv.Itoa(status)).Observe(float64(duration))
		}()
		ctx.Next() //先继续往下执行，最后在执行这个func，计算总的http相应时间
	}
}
func (b *Builder) BuildActiveRequest() gin.HandlerFunc {
	gauge := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: b.Namespace,
		Subsystem: b.Subsystem,
		Help:      b.Help,
		Name:      b.Name + "_active_req",
		ConstLabels: map[string]string{
			"instance_id": b.InstanceId,
		},
	})
	prometheus.MustRegister(gauge)
	return func(ctx *gin.Context) {
		gauge.Inc()
		defer gauge.Dec()
		ctx.Next()
	}

}
