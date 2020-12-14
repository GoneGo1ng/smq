package main

import (
	"github.com/Allenxuxu/gev"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"hzsun.com/easyconnect/smq/broker"
)

func main() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	cfg.OutputPaths = []string{"stdout"}
	logger, err := cfg.Build()
	if err != nil {
		panic(err)
	}
	zap.ReplaceGlobals(logger)

	handler, err := broker.LoadBroker()
	if err != nil {
		panic(err)
	}

	s, err := gev.NewServer(handler, gev.Address(":1883"))
	if err != nil {
		panic(err)
	}
	zap.L().Info("启动成功")
	s.Start()
}
