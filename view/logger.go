package view

import (
	"vincentinttsh/openclass-system/internal/mode"

	"go.uber.org/zap"
)

// Logger is a global logger
var Logger *zap.Logger
var sugar *zap.SugaredLogger

const (
	infoMsgLevel = "info"
	warnMsgLevel = "warn"
	errMsgLevel  = "error"
)

func initLogger() {
	if mode.Mode() == mode.ReleaseMode {
		Logger, _ = zap.NewProduction()
	} else if mode.Mode() == mode.DebugMode {
		Logger, _ = zap.NewDevelopment()
	}
	sugar = Logger.Sugar()
}
