package main

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// 支持静态令牌鉴权 - 目前仅支持 env
func authMiddleware() gin.HandlerFunc {
	const defaultRealm = "Coffee Time"

	var accounts = make(map[string]string)
	for _, v := range os.Environ() {
		tmp := strings.SplitN(v, "=", 2)
		user, pw := tmp[0], tmp[1]
		if !strings.HasPrefix(user, "user_") {
			continue
		}
		accounts[strings.TrimPrefix(user, "user_")] = pw
	}

	if len(accounts) == 0 {
		return nil
	}

	zap.L().Info("basic auth is enable")
	return gin.BasicAuthForRealm(accounts, defaultRealm)
}
