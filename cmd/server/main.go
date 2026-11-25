package main

import (
	"context"
	"log/slog"
	"monitoring/internal/logger"
	"monitoring/internal/server"
	"monitoring/internal/storage"
	metricsv1 "monitoring/pkg/api/metrics/v1"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	logger := logger.New(os.Stdout, slog.LevelDebug)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	//dbStorage := storage.NewStorageEmulator(logger)
	dbStorage, err := storage.New("database/monitoring.db")
	if err != nil {
		logger.Error("failed to init db", "err", err)
		return
	}
	defer dbStorage.Close() // Обязательно закрываем дескриптор файла при выходе

	latest, err := dbStorage.GetLatest(ctx, 3)
	if err != nil {
		logger.Error("test read failed", "err", err)
	} else {
		logger.Info("test read success", "count", len(latest), "data", latest)
	}

	scraper, _ := server.NewScraper("http://localhost:8080", dbStorage, logger)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		scraper.Run(ctx)
	}()

	// 1. Открываем TCP-порт для gRPC
	grpcLis, err := net.Listen("tcp", ":50051")
	if err != nil {
		logger.Error("failed to listen grpc port", "err", err)
		return
	}

	// 2. Создаем сам gRPC-сервер
	grpcServer := grpc.NewServer()

	// 3. Регистрируем наш сервис и включаем рефлексию для grpcurl
	metricsService := server.NewGRPCServer(dbStorage)
	metricsv1.RegisterMetricsServiceServer(grpcServer, metricsService)
	reflection.Register(grpcServer) // <--- ВОТ ОНА, МАГИЯ ДЛЯ ОТЛАДКИ!

	// 4. Запускаем gRPC в отдельной горутине
	wg.Add(1)
	go func() {
		defer wg.Done()
		logger.Info("starting grpc server", "addr", ":50051")
		if err := grpcServer.Serve(grpcLis); err != nil {
			logger.Error("grpc server failed", "err", err)
		}
	}()

	<-ctx.Done()
	wg.Wait()

	// 5. В секцию Graceful Shutdown (после <-ctx.Done()):
	grpcServer.GracefulStop() // вежливо ждет завершения активных RPC

	logger.Info("server stopped gracefully")
}
