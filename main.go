package main

import (
	"math/rand"
	"os"
	"strings"
	"time"
	_ "time/tzdata"

	_ "bark-serverless/controller"
	"bark-serverless/logger"
	"bark-serverless/router"

	apiWrapper "github.com/TMaize/scf-apigw-wrap"
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
	ApiServerName = "api_path" // API网关名称
	runModeKey    = "mode"     // 运行环境
	portKey       = "port"     // API服务运行时的端口
)

func init() {
	// 初始化时区
	time.Local, _ = time.LoadLocation("Asia/Shanghai")

	// 初始化随机数种子
	rand.Seed(time.Now().UnixNano())
}

func main() {
	switch mode := runMode(cast.ToUint(os.Getenv(runModeKey))); mode {
	case runModeTencentSCF:
		logger.GetGlobalLog().Info("SCF模式")

		cloudfunction.Start(func(req events.APIGatewayRequest) (resp events.APIGatewayResponse, err error) {
			// 转换一下路由路径
			path := strings.TrimPrefix(req.Path, os.Getenv(ApiServerName))
			if !strings.HasPrefix(path, "/") {
				path = "/" + path
			}

			// 调用gin server
			resp = apiWrapper.Wrap(req, path, router.SetupRouter())

			return
		})
	case runModeApiServer:
		logger.GetGlobalLog().Info("api服务模式")

		if err := router.SetupRouter().Run(":" + func() string {
			p := os.Getenv(portKey)
			if p == "" {
				return "8080"
			}
			return p
		}()); err != nil {
			logger.GetGlobalLog().Error("服务启动失败", zap.Error(err))
			return
		}
	}
	return
}
