package controller

import (
	"net/http"
	"time"

	"bark-serverless/router"

	"github.com/gin-gonic/gin"
)

func init() {
	router.AppendRouter(func(r *gin.Engine) {
		r.GET("/ping", Ping)
		r.GET("/healthz", Healthz)
	})
}

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, CommonResp{
		Code:      http.StatusOK,
		Message:   "pong",
		Timestamp: time.Now().Unix(),
		Data:      nil,
	})
}

func Healthz(c *gin.Context) {
	c.String(http.StatusOK, "ok")
}
