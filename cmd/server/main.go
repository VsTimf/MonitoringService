package main

import (
	"context"
	"log/slog"
	"monitoring/internal/collector"
	"monitoring/internal/logger"
	"monitoring/internal/storage"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	logger := logger.New(os.Stdout, slog.LevelDebug)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	//dbStorage := storage.NewStorageEmulator(logger)
	dbStorage, err := storage.New("metrics.db")
	if err != nil {
		logger.Error("failed to init db", "err", err)
		return
	}
	defer dbStorage.Close() // Обязательно закрываем дескриптор файла при выходе

	scraper, _ := collector.NewScraper("http://localhost:8080", dbStorage, logger)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		scraper.Run(ctx)
	}()

	<-ctx.Done()
	wg.Wait()
	logger.Info("server stopped gracefully")
}
