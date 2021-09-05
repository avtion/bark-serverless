package main

import "github.com/gin-gonic/gin"

func setupDefaultRouter(r gin.IRoutes) {
	// 杂项
	r.GET("/ping", Ping)
	r.GET("/healthz", Healthz)

	// 推送 v2
	r.POST("/push", Push)

	// 推送 v1
	r.GET("/:device_key/:title", Push)
	r.POST("/:device_key/:title", Push)
	r.GET("/:device_key/:title/*body", Push)
	r.POST("/:device_key/:title/*body", Push)

	// 设备注册
	r.POST("/register", RegisterDevice)
	r.GET("/register", RegisterDevice)
}