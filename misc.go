package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

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
