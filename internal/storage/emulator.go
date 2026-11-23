package storage

import (
	"context"
	"log/slog"
	"monitoring/internal/domain"
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
