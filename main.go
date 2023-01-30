package main

import (
	"github.com/gin-gonic/gin"

	"github.com/lwangrabbit/gin-http-metrics/middleware"
)

func main() {
	r := gin.New()

	r.Use(middleware.HttpMetricMiddleware())
	r.GET("/metrics", middleware.HttpMetricHandler())

	if err := r.Run(":8081"); err != nil {
		panic(err)
	}
}
