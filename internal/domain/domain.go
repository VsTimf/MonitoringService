package domain

import (
	"fmt"
	"time"
)

const (
	_   = 1 << (10 * iota) // 0: пропускаем (1 << 0 = 1 байт)
	KiB                    // 1 << 10 (1024)
	MiB                    // 1 << 20 (1048576)
	GiB                    // 1 << 30 (1073741824)
)

type Metrics struct {
	Host      string    `json:"host"`
	CPU       float64   `json:"cpu"`
	RAM       uint64    `json:"ram"`
	Timestamp time.Time `json:"timestamp"`
}

func (m Metrics) String() string {
	return fmt.Sprintf("Host: %s CPU: %.1f RAM:%.1fGb Time:%v", m.Host, m.CPU, float64(m.RAM)/GiB, m.Timestamp.Format(time.TimeOnly))
}
