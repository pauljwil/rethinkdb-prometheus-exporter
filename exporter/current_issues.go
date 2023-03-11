package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	gorethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type issue struct {
	ID          string `rethinkdb:"id"`
	Type        string `rethinkdb:"type"`
	Critical    bool   `rethinkdb:"critical"`
	Description string `rethinkdb:"description"`
}

type issueTotals struct {
	logWriteError      float64
	nameCollision      float64
	outdatedIndex      float64
	tableAvailability  float64
	memoryError        float64
	nonTransitiveError float64
	other              float64
}

func (e *RethinkdbExporter) collectRethinkCurrentIssues(ch chan<- prometheus.Metric) int {
	errcount := 0

	cur, err := gorethink.DB(gorethink.SystemDatabase).Table(gorethink.CurrentIssuesSystemTable).Run(e.rconn)
	if err != nil {
		log.Error().Err(err).Msg("failed to query system current_issues table")
		errcount++
		return errcount
	}
	defer func() {
		err := cur.Close()
		if err != nil {
			log.Warn().Err(err).Msg("error while closing cursor")
		}
	}()

	if cur.Err() != nil {
		log.Error().Err(cur.Err()).Msg("query error from cursor")
		errcount++
		return errcount
	}

	var (
		ci  issue
		cit issueTotals
	)
	for cur.Next(&ci) {
		if cur.Err() != nil {
			log.Error().Err(cur.Err()).Msg("query error from cursor")
			errcount++
			return errcount
		}

		countCurrentIssue(ci, &cit)
	}

	e.processCurrentIssues(cit, ch)

	return errcount
}

func countCurrentIssue(ci issue, cit *issueTotals) {
	switch ci.Type {
	case "log_write_error":
		cit.logWriteError++
	case "server_name_collision":
		cit.nameCollision++
	case "db_name_collision":
		cit.nameCollision++
	case "table_name_collision":
		cit.nameCollision++
	case "outdated_index":
		cit.outdatedIndex++
	case "table_availability":
		cit.tableAvailability++
	case "memory_error":
		cit.memoryError++
	case "non_transitive_error":
		cit.nonTransitiveError++
	default:
		cit.other++
	}
}

func (e *RethinkdbExporter) processCurrentIssues(cit issueTotals, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(e.metrics.logWriteIssues, prometheus.GaugeValue, cit.logWriteError)
	ch <- prometheus.MustNewConstMetric(e.metrics.nameCollisionIssues, prometheus.GaugeValue, cit.nameCollision)
	ch <- prometheus.MustNewConstMetric(e.metrics.outdatedIndexIssues, prometheus.GaugeValue, cit.outdatedIndex)
	ch <- prometheus.MustNewConstMetric(e.metrics.totalAvailabilityIssues, prometheus.GaugeValue, cit.tableAvailability)
	ch <- prometheus.MustNewConstMetric(e.metrics.memoryAvailabilityIssues, prometheus.GaugeValue, cit.memoryError)
	ch <- prometheus.MustNewConstMetric(e.metrics.connectivityIssues, prometheus.GaugeValue, cit.nonTransitiveError)
	ch <- prometheus.MustNewConstMetric(e.metrics.otherIssues, prometheus.GaugeValue, cit.other)
}
