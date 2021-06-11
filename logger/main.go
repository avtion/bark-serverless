package logger

import "go.uber.org/zap"

var log *zap.Logger

func init() {
	SetupLogger()
}

func SetupLogger() {
	log, _ = zap.NewProduction(zap.AddCaller())
}

func GetGlobalLog() *zap.Logger {
	if log == nil {
		SetupLogger()
	}
	return log
}
