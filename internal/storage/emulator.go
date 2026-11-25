package storage

import (
	"context"
	"log/slog"
	"monitoring/internal/domain"
	"time"
)

type StorageEmulator struct {
	logger *slog.Logger
}

func NewStorageEmulator(logger *slog.Logger) *StorageEmulator {
	if logger == nil {
		logger = slog.Default()
	}

	return &StorageEmulator{
		logger: logger.With("m", "storage emulator"),
	}
}

func (s *StorageEmulator) Save(ctx context.Context, m domain.Metrics) error {
	s.logger.Info("Saved", "metrics", m)
	return nil
}

func (s *StorageEmulator) GetLatest(ctx context.Context, limit int) ([]domain.Metrics, error) {

	metrics := []domain.Metrics{
		{Host: "fedora", CPU: 30.0, RAM: 1073741824, Timestamp: time.Now()},
		{Host: "fedora", CPU: 31.0, RAM: 1073741824, Timestamp: time.Now()},
		{Host: "fedora", CPU: 32.0, RAM: 1073741824, Timestamp: time.Now()},
	}

	return metrics, nil
}
