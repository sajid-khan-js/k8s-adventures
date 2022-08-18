package main

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

/*

https://prometheus.io/docs/guides/go-application/
https://gabrieltanner.org/blog/collecting-prometheus-metrics-in-golang/

*/

var (
	totalRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Number of HTTP requests.",
		},
		[]string{"method", "path"})

	responseStatus = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "response_status",
			Help: "Status of HTTP responses.",
		},
		[]string{"method", "path", "status"})

	httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "http_response_time_seconds",
		Help: "Duration of HTTP requests.",
	}, []string{"method", "path"})
)

func prometheusMiddleware(c *gin.Context) {

	// Pre-handlers
	method := c.Request.Method
	path := c.Request.URL.Path

	// Prometheus itself will let us know if there are issues with /metrics
	if path == "/metrics" {
		c.Next()
		return
	}

	// Ensure we don't generate time series for every possible namespace so /namespaces/kube-system becomes /namespaces/:name
	for _, p := range c.Params {
		if p.Key == "name" {
			path = strings.Replace(path, p.Value, ":name", 1)
			break
		}
	}

	totalRequests.WithLabelValues(method, path).Inc()
	timer := prometheus.NewTimer(httpDuration.WithLabelValues(method, path))

	// Pass on to the next-in-chain i.e. the handlers
	c.Next()

	// Post-handlers
	statusCode := strconv.Itoa(c.Writer.Status())
	responseStatus.WithLabelValues(method, path, statusCode).Inc()
	timer.ObserveDuration()

}

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
}
