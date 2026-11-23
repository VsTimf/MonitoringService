package main

import (
	"context"
	"log/slog"
	"monitoring/internal/agent"
	"monitoring/internal/logger"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func main() {

	logger := logger.New(os.Stdout, slog.LevelDebug)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	measurer := agent.NewMeasurer(logger)
	server := agent.NewServer(":8080", measurer, logger)

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		measurer.Run(ctx)
	}()

	go func() {
		if err := server.Start(); err != nil {
			logger.Error("server failed", "err", err)
			stop() // если порт занят — роняем все приложение
		}
	}()

	<-ctx.Done()

	// Даем серверу максимум 3 секунды на закрытие коннектов
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := server.Stop(shutdownCtx); err != nil {
		logger.Error("failed to stop server gracefully", "err", err)
	}

	logger.Info("agent stopped cleanly")

	wg.Wait()

}
