package storage

import (
	"context"
	"database/sql"
	"fmt"
	"monitoring/internal/domain"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

type SQLiteStorage struct {
	db *sql.DB
}

func New(dbPath string) (*SQLiteStorage, error) {

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("create db directory %s: %w", dir, err)
	}

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

func (s *SQLiteStorage) GetLatest(ctx context.Context, limit int) ([]domain.Metrics, error) {
	query := `
	SELECT host, cpu_usage, ram_usage, timestamp
	FROM metrics
	ORDER BY id DESC
	LIMIT ?;
	`

	// Запрос
	rows, err := s.db.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("query metrics: %w", err)
	}
	defer rows.Close()

	// Предаллоцируем срез под нужный размер для экономии памяти
	result := make([]domain.Metrics, 0, limit)

	for rows.Next() {
		var (
			m      domain.Metrics
			tsUnix int64
		)

		if err := rows.Scan(&m.Host, &m.CPU, &m.RAM, &tsUnix); err != nil {
			return nil, fmt.Errorf("scan metric row: %w", err)
		}

		m.Timestamp = time.Unix(tsUnix, 0)
		result = append(result, m)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration: %w, err")
	}

	return result, nil

}
