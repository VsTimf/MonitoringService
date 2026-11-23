package collector

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"monitoring/internal/domain"
	"net/http"
	"net/url"
	"time"
)

type Storage interface {
	Save(ctx context.Context, m domain.Metrics) error

	//GetLatest(ctx context.Context, host string, limit int) ([]domain.Metrics, error)
}

const (
	polltime = 2 * time.Second
)

type Scraper struct {
	logger   *slog.Logger
	agentURL string
	client   *http.Client
	storage  Storage
}

func NewScraper(urlStr string, storage Storage, logger *slog.Logger) (*Scraper, error) {
	if logger == nil {
		logger = slog.Default()
	}

	//check url spelling
	parsedURL, err := url.ParseRequestURI(urlStr)

	if err != nil {
		return nil, fmt.Errorf("invalid agent url: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return nil, fmt.Errorf("unsupported protocol scheme: %s", parsedURL.Scheme)
	}

	//create http client
	client := &http.Client{
		Timeout: 2 * time.Second,
	}

	return &Scraper{
		logger:   logger.With("m", "server"),
		agentURL: urlStr,
		client:   client,
		storage:  storage,
	}, nil
}

func (s *Scraper) Run(ctx context.Context) {
	ticker := time.NewTicker(polltime)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return

		case <-ticker.C:
			s.fetch(ctx)
		}
	}
}

func (s *Scraper) fetch(ctx context.Context) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.agentURL+"/metrics", nil)
	if err != nil {
		s.logger.Error("failed to create request", "err", err)
		return
	}

	resp, err := s.client.Do(req)
	if err != nil {
		s.logger.Warn("agent unreachable", "err", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		s.logger.Error("unexprcted status code", "status", resp.StatusCode)
		return
	}

	var m domain.Metrics
	if err := json.NewDecoder(resp.Body).Decode(&m); err != nil {
		s.logger.Error("failed to decode json", "err", err)
		return
	}

	if err := s.storage.Save(ctx, m); err != nil {
		s.logger.Error("failed to save metrics", "err", err)
	}
}
