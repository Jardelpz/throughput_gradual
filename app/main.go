package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"time"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"status", "method"},
	)

	httpRequestDuration = prometheus.NewSummaryVec(
		prometheus.SummaryOpts{
			Name: "http_request_duration_milliseconds",
			Help: "HTTP request latencies in milliseconds",
		},
		[]string{"status", "method"},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
}

func main() {
	r := gin.Default()
	r.GET("/hc", func(c *gin.Context) {
		statusCode := http.StatusOK
		start := time.Now()

		c.JSON(http.StatusOK, gin.H{
			"message": "alive and kicking",
		})

		elapsed := float64(time.Since(start).Milliseconds())
		httpRequestsTotal.WithLabelValues(strconv.Itoa(statusCode), c.Request.Method).Inc()
		httpRequestDuration.WithLabelValues(strconv.Itoa(statusCode), c.Request.Method).Observe(elapsed)

	})

	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	r.Run()
}
