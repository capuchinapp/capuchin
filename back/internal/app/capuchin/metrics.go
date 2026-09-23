package capuchin

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"capuchin/internal/metrics"
)

const metricsReadHeaderTimeout = 5 * time.Second

// startMetrics запускает сервер метрик и возвращает его для graceful shutdown.
func startMetrics(
	port int,
	db *sql.DB,
	dbName string,
	metricsInstance *metrics.Metrics,
	logger *zap.Logger,
) *http.Server {
	reg := prometheus.NewRegistry()

	reg.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewDBStatsCollector(db, dbName),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		metricsInstance.HTTPRequestsTotal,
		metricsInstance.HTTPRequestsCurrent,
		metricsInstance.HTTPRequestDurationSeconds,
	)

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	server := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		Handler:           mux,
		ReadHeaderTimeout: metricsReadHeaderTimeout,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("metrics listen",
				zap.Error(err),
			)
			panic(err)
		}
	}()

	logger.Info("Metrics started",
		zap.String("address", server.Addr),
	)

	return server
}
