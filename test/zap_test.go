package main

import (
	"go.uber.org/zap"
	"grpc-lb/internal/pkg/log"
	"testing"
	"time"
)

func TestUploadConfig(t *testing.T) {
	url := "Hello"
	//logger, _ := zap.NewProduction()
	logger, _ := zap.NewDevelopment()

	//Sync刷新任何缓冲的日志条目。
	defer logger.Sync()
	logger.Info("failed to fetch URL",
		// Structured context as strongly typed Field values.
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second),
	)
	logger.Warn("debug log", zap.String("level", url))
	logger.Error("Error Message", zap.String("error", url))
}

func TestGetLogger(t *testing.T) {
	url := "test"
	logger := log.GetLogger()
	logger.Info("test for info", zap.String("url", url))
}

func BenchmarkUploadConfig2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		url := "test"
		logger := log.GetLogger()
		logger.Info("test for info", zap.String("url", url))
	}
}
