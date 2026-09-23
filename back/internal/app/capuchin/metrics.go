package capuchin

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"

	"capuchin/internal/metrics"
)

// startMetrics запускает метрики.
func startMetrics(
	port int,
	db *sql.DB,
	dbName string,
	metricsInstance *metrics.Metrics,
	logger *zap.Logger,
) {
	metricsAddr := fmt.Sprintf(":%d", port)

	go func() {
		reg := prometheus.NewRegistry()

		reg.MustRegister(
			collectors.NewGoCollector(),
			collectors.NewDBStatsCollector(db, dbName),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
			metricsInstance.HTTPRequestsTotal,
			metricsInstance.HTTPRequestsCurrent,
			metricsInstance.HTTPRequestDurationSeconds,
		)

		http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

		if err := http.ListenAndServe(metricsAddr, nil); err != nil { //nolint:gosec // Данный адрес не будет доступен в интернете
			logger.Error("metrics listen",
				zap.Error(err),
			)
			panic(err)
		}
	}()

	logger.Info("Metrics started",
		zap.String("address", metricsAddr),
	)
}
