package internal_test

import (
	"testing"
	"time"

	"github.com/fishykins/aero/internal"
	"go.uber.org/zap"
)

func TestApp(t *testing.T) {
	now := time.Now()
	logger := zap.NewNop()
	app := internal.App{
		Logger: logger.Sugar(),
	}
	app.Init()
	app.Run(&now)
	app.Shutdown()
}
