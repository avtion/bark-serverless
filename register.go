package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.uber.org/zap"
)

// RegisterDevice 注册设备
func RegisterDevice(c *gin.Context) {
	// 实现Bark-Server的注册路由
	// 将会打印device_key到控制台，开发者需要手动将device_key配置到Serverless全局环境变量中

	// 初始化日志
	l := zap.L().With(zap.String("router", "register"))

	// 初始化请求
	var (
		req = new(DeviceInfo)
		err error
	)

	// 先尝试从GET请求中获取device_key
	if err = c.ShouldBindQuery(req); err != nil {
		l.Error("设备注册失败", zap.Error(err))

		c.AbortWithStatusJSON(http.StatusInternalServerError, CommonResp{
			Code:      http.StatusInternalServerError,
			Message:   "device registration failed",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	// 如果有DeviceToken则跳到响应
	if req.OldDeviceToken != "" {
		req.DeviceKey, req.DeviceToken = req.OldDeviceKey, req.OldDeviceToken
	}

	// 如果GET请求参数解析不到就从Body获取
	if (req.DeviceKey == "" || req.DeviceToken == "") && c.Request.Method == http.MethodPost {
		if err = c.ShouldBindBodyWith(req, binding.JSON); err != nil {
			l.Error("设备注册失败", zap.Error(err))

			c.AbortWithStatusJSON(http.StatusInternalServerError, CommonResp{
				Code:      http.StatusInternalServerError,
				Message:   "device registration failed",
				Timestamp: time.Now().Unix(),
			})
			return
		}
	}

	if req.DeviceToken == "" || req.OldDeviceToken == "" {
		l.Error("设备注册失败，token为空")

		c.AbortWithStatusJSON(http.StatusBadRequest, CommonResp{
			Code:      http.StatusBadRequest,
			Message:   "device token is empty",
			Timestamp: time.Now().Unix(),
		})
		return
	}

	l = l.With(
		zap.String("key", req.DeviceKey), zap.String("token", req.DeviceToken),
		zap.String("old_key", req.OldDeviceKey), zap.String("old_token", req.OldDeviceToken),
	)
	l.Info("设备绑定信息")

	c.JSON(http.StatusOK, CommonResp{
		Code:    http.StatusOK,
		Message: "success",
		Data: map[string]string{
			"key":          req.DeviceKey,
			"device_key":   req.DeviceKey,
			"device_token": req.DeviceToken,
		},
		Timestamp: time.Now().Unix(),
	})

	// 写入数据库不影响核心业务
	for _, db := range readDBFromCtx(c) {
		if v, ok := db.(SaveToken); ok {
			if err = v.SaveToken(req.DeviceKey, req.DeviceToken); err != nil {
				l.Error("save token failed", zap.Error(err))
				continue
			}
			l.Info("save token successfully", zap.String("dbName", db.Name()))
		}
	}
}
