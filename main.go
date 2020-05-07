package main

import (
	"context"

	"github.com/etiennemarais/world-cities-database/cmd/list"
	"github.com/etiennemarais/world-cities-database/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func main() {
	ctx := context.Background()
	logger := logger.New()
	logger.Info("booting")

	defer func() {
		_ = logger.Sync()
	}()

	var cmd = &cobra.Command{Use: "world-cities-database"}

	cmd.AddCommand(
		list.Command(ctx, logger),
	)

	if err := cmd.ExecuteContext(ctx); err != nil {
		logger.Fatal("execute command", zap.Error(err))
	}
}
