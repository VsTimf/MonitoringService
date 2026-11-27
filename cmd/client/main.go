package main

import (
	"context"
	"log/slog"
	"monitoring/internal/client"
	"monitoring/internal/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	logger := logger.New(os.Stdout, slog.LevelDebug)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup

	client := client.New(logger)

	wg.Add(1)
	go func() {
		defer wg.Done()
		client.Run(ctx)
	}()

	<-ctx.Done()
	wg.Wait()
}
