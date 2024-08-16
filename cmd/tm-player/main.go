package main

import (
	"context"

	"github.com/javascriptizer1/tm-player.backend/internal/app"
	"github.com/javascriptizer1/tm-shared.backend/logger"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()

	a, err := app.New(ctx)
	if err != nil {
		logger.Fatal("failed to init app ", zap.Error(err))
	}

	if err := a.Run(); err != nil {
		logger.Fatal("failed to run app: ", zap.Error(err))
	}
}
