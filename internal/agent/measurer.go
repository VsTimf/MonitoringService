package agent

import (
	"context"
	"log/slog"
	"monitoring/internal/domain"
	"os"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

const (
	polltime = 1 * time.Second
)

type Measurer struct {
	logger *slog.Logger

	mu      sync.RWMutex
	metrics domain.Metrics
}

func NewMeasurer(logger *slog.Logger) *Measurer {

	if logger == nil {
		logger = slog.Default()
	}

	hostname, err := os.Hostname()

	if err != nil {
		hostname = "Unknown host"
	}

	return &Measurer{
		logger:  logger.With("m", "Measurer"),
		metrics: domain.Metrics{Host: hostname},
	}
}

func (m *Measurer) Run(ctx context.Context) {

	ticker := time.NewTicker(polltime)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			m.measure()
		}
	}
}

func (m *Measurer) measure() {

	// Получить значение CPU
	cpuPercentages, err := cpu.Percent(0, false)
	if err != nil {
		m.logger.Error("failed to get cpu stats", "err", err)
		return
	}

	var cpuVal float64
	if len(cpuPercentages) > 0 {
		cpuVal = cpuPercentages[0]
	}

	// Получить значение RAM
	vMem, err := mem.VirtualMemory()
	if err != nil {
		m.logger.Error("failed to get memory stats", "err", err)
		return
	}

	now := time.Now()

	// Сохранить под мьютексом
	m.mu.Lock()
	m.metrics.CPU = cpuVal
	m.metrics.RAM = vMem.Used
	m.metrics.Timestamp = now
	current := m.metrics // копируем структуру для безопасного логирования
	m.mu.Unlock()

	m.logger.Info("metric collected", "metrics", current)
}

func (m *Measurer) Metrics() domain.Metrics {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.metrics
}
