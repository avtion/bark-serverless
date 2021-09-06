package main

import (
	"math/rand"
	"os"
	"strings"
	"time"
	_ "time/tzdata"

	apiWrapper "github.com/TMaize/scf-apigw-wrap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"github.com/tencentyun/scf-go-lib/cloudfunction"
	"github.com/tencentyun/scf-go-lib/events"
	"go.uber.org/zap"
)

// 程序运行模式
type runMode uint

const (
	runModeTencentSCF runMode = iota // 腾讯云SCF - 用于SCF部署
	runModeApiServer                 // API服务 - 用于本地调试
)

var _ = [...]runMode{runModeTencentSCF, runModeApiServer}

const (
	apiServerName = "api_path" // API网关名称
	runModeKey    = "mode"     // 运行环境
	portKey       = "port"     // API服务运行时的端口
)

func main() {
	var log = initLogger()

	// 初始化时区
	time.Local, _ = time.LoadLocation("Asia/Shanghai")

	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())

	// Gin
	engine := gin.New()
	setupDefaultRouter(engine.Use(WriteTokenStoreToCtx(store...)))

	switch mode := runMode(cast.ToUint(os.Getenv(runModeKey))); mode {
	case runModeTencentSCF:
		log.Info("SCF 模式")

		cloudfunction.Start(func(req events.APIGatewayRequest) (resp events.APIGatewayResponse, err error) {
			// 转换一下路由路径
			path := strings.TrimPrefix(req.Path, os.Getenv(apiServerName))
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}

			// 调用gin server
			resp = apiWrapper.Wrap(req, path, engine)

			return
		})
	case runModeApiServer:
		log.Info("API 服务模式")

		if err := engine.Run(":" + func() string {
			p := os.Getenv(portKey)
			if p == "" {
				return "8080"
			}
			return p
		}()); err != nil {
			log.Desugar().Error("服务启动失败", zap.Error(err))
			return
		}
	}
}
