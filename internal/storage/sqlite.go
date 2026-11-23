package storage

import (
	"context"
	"database/sql"
	"fmt"
	"monitoring/internal/domain"
	"time"

	_ "modernc.org/sqlite"
)

/*
type Storage interface {
    Save(ctx context.Context, m domain.Metrics) error
}
*/

type SQLiteStorage struct {
	db *sql.DB
}

func New(dbPath string) (*SQLiteStorage, error) {

	dsn := fmt.Sprintf("%s?_journal_mode=WAL&_busy_timeout=5000", dbPath)

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping db failed: %w", err)
	}

	query := `
	CREATE TABLE IF NOT EXISTS metrics (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	host TEXT NOT NULL,
	cpu_usage REAL NOT NULL,
	ram_usage INTEGER NOT NULL,
	timestamp INTEGER NOT NULL
	);`

	if _, err := db.ExecContext(ctx, query); err != nil {
		return nil, fmt.Errorf("create table failed: %w", err)
	}

	return &SQLiteStorage{db: db}, nil
}

func (s *SQLiteStorage) Close() error {
	return s.db.Close()
}

func (s SQLiteStorage) Save(ctx context.Context, m domain.Metrics) error {
	query := `
	INSERT INTO metrics (host, cpu_usage, ram_usage, timestamp)
	VALUES (?, ?, ?, ?);
	`

	_, err := s.db.ExecContext(
		ctx,
		query,
		m.Host,
		m.CPU,
		m.RAM,
		m.Timestamp.Unix(),
	)
	if err != nil {
		return fmt.Errorf("insert metric: %w", err)
	}

	return nil
}
