package agent

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"monitoring/internal/domain"
	"net/http"
	"time"
)

type MetricsProvider interface {
	Metrics() domain.Metrics
}

type Server struct {
	logger     *slog.Logger
	httpServer *http.Server
}

func NewServer(addr string, provider MetricsProvider, logger *slog.Logger) *Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", handleMetrics(provider))

	srv := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 2 * time.Second,
		WriteTimeout:      5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	if logger == nil {
		logger = slog.Default()
	}

	return &Server{
		httpServer: srv,
		logger:     logger.With("m", "httpserver"),
	}
}

func (s *Server) Start() error {
	s.logger.Info("Starting http server", "addr", s.httpServer.Addr)

	err := s.httpServer.ListenAndServe()

	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("stopping http server")
	return s.httpServer.Shutdown(ctx)
}

func handleMetrics(provider MetricsProvider) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		data := provider.Metrics()

		if err := json.NewEncoder(w).Encode(data); err != nil {
			http.Error(w, "failed to encode json", http.StatusInternalServerError)
		}
	}
}
