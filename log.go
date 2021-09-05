package main

import "go.uber.org/zap"

func initLogger() *zap.SugaredLogger {
	customLogger, _ := zap.NewProduction()
	zap.ReplaceGlobals(customLogger)
	return customLogger.Sugar()
}
