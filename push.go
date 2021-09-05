package main

import (
	"net/http"
	"strings"
	"time"

	"github.com/finb/bark-server/v2/apns"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

// Push 推送
func Push(c *gin.Context) {
	// 初始化日志
	l := zap.L().With(zap.String("router", "push"))

	// 初始化参数
	var req = &apns.PushMessage{
		DeviceToken: "",
		DeviceKey:   "",
		Category:    "myNotificationCategory",
		Title:       "",
		Body:        "NoContent",
		Sound:       "1107",
		ExtParams:   make(map[string]interface{}),
	}

	// 首先尝试从请求内容体从获取参数
	// https://blog.csdn.net/yes169yes123/article/details/106204252
	if c.Request.Method == http.MethodPost {
		if err := c.ShouldBindBodyWith(req, binding.JSON); err != nil {
			l.Error("绑定数据失败", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusBadRequest, CommonResp{
				Code:      http.StatusBadRequest,
				Message:   "request bind failed",
				Timestamp: time.Now().Unix(),
			})
			return
		}
	}

	// 读取query参数并不会触发内容体Reader的EOF
	// query的参数会覆盖掉Body的参数
	if err := c.ShouldBindQuery(req); err != nil {
		l.Error("绑定数据失败", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusBadRequest, CommonResp{
			Code:      http.StatusBadRequest,
			Message:   "request bind failed",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// 再尝试从URI从获取参数
	deviceKey, isExist := c.Params.Get("device_key")
	if isExist {
		req.DeviceKey = deviceKey
	}
	title, isExist := c.Params.Get("title")
	if isExist {
		req.Title = title
	}
	body, isExist := c.Params.Get("body")
	if isExist {
		req.Body = strings.TrimPrefix(body, "/")
	}

	// 如果设备的Key为空则中断流程
	if req.DeviceKey == "" {
		l.Error("device key is empty")
		c.AbortWithStatusJSON(http.StatusBadRequest, CommonResp{
			Code:      http.StatusBadRequest,
			Message:   "device key is empty",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	dbs := readDBFromCtx(c)
	for _, db := range dbs { // 根据range特性, 如果dbs为nil会直接跳过循环
		if v, ok := db.(LoadToken); ok {
			req.DeviceToken, _ = v.LoadToken(req.DeviceKey)
		}
	}
	if req.DeviceToken == "" {
		l.Error("failed to get token from db", zap.String("key", req.DeviceKey))
		c.AbortWithStatusJSON(http.StatusBadRequest, CommonResp{
			Code:      http.StatusBadRequest,
			Message:   "failed to get token from db",
			Timestamp: time.Now().Unix(),
		})
	}

	// 推送消息
	if err := apns.Push(req); err != nil {
		l.Error("failed to push message", zap.Error(err))
		c.AbortWithStatusJSON(http.StatusInternalServerError, CommonResp{
			Code:      http.StatusInternalServerError,
			Message:   "failed to push message",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	c.JSON(http.StatusOK, CommonResp{
		Code:      http.StatusOK,
		Message:   "success",
		Timestamp: time.Now().Unix(),
	})
}
