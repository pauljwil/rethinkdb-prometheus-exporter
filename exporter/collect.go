package exporter

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
)

// Collect sends collected metrics values to the prometheus chan
func (e *RethinkdbExporter) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()

	ctx := context.TODO() // TODO: add scrape timeout
	errcount := e.collectRethinkStats(ctx, ch)
	errcount = errcount + e.collectRethinkCurrentIssues(ch)
	errcount = errcount + e.collectRethinkTableSizes(ch)

	elapsed := time.Since(start)
	ch <- prometheus.MustNewConstMetric(e.metrics.scrapeErrors, prometheus.GaugeValue, float64(errcount))
	ch <- prometheus.MustNewConstMetric(e.metrics.scrapeLatency, prometheus.GaugeValue, elapsed.Seconds())

	log.Debug().Dur("duration", elapsed).Msg("collect finished")
}
