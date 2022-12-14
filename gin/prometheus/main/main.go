package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"math/rand"
	"net/http"
	"net/url"
	"strconv"
)

func main() {
	engine := gin.New()

	engine.GET("/hello", func(context *gin.Context) {
		context.String(http.StatusOK, "Hello")
	})

	engine.GET("/counter", func(context *gin.Context) {
		purl, _ := url.Parse(context.Request.RequestURI)
		AccessCounter.With(prometheus.Labels{
			"method": context.Request.Method,
			"path":   purl.Path,
		}).Add(1)
	})

	engine.GET("queue", func(context *gin.Context) {
		num := context.Query("num")
		fnum, _ := strconv.ParseFloat(num, 32)
		QueueGauge.With(prometheus.Labels{"name": "queue_gevin"}).Set(fnum)
	})

	engine.GET("/histogram", func(context *gin.Context) {
		purl, _ := url.Parse(context.Request.RequestURI)
		HttpDurationsHistogram.With(prometheus.Labels{"path": purl.Path}).Observe(float64(rand.Intn(30)))
	})

	engine.GET("/summary", func(c *gin.Context) {
		purl, _ := url.Parse(c.Request.RequestURI)
		HttpDurations.With(prometheus.Labels{"path": purl.Path}).Observe(float64(rand.Intn(30)))
	})

	engine.GET("/metrics", gin.WrapH(promhttp.Handler()))

	_ = engine.Run(":8000")
}
