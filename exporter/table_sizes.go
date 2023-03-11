package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog/log"
	gorethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

type tableSize struct {
	DB    string  `rethinkdb:"db"`
	Size  float64 `rethinkdb:"size"`
	Table string  `rethinkdb:"table"`
}

func (e *RethinkdbExporter) collectRethinkTableSizes(ch chan<- prometheus.Metric) int {
	errcount := 0

	cur, err := gorethink.DB(gorethink.SystemDatabase).
		Table(gorethink.StatsSystemTable).
		HasFields("db", "table").
		Group("db", "table").
		Map(func(doc gorethink.Term) gorethink.Term {
			return doc.Field("storage_engine").Field("disk").Field("space_usage").Field("data_bytes").Default(0)
		}).
		Sum().
		Ungroup().
		Map(func(doc gorethink.Term) interface{} {
			return map[string]interface{}{
				"db":    doc.Field("group").Nth(0),
				"table": doc.Field("group").Nth(1),
				"size":  doc.Field("reduction").Div(1024).Div(1024),
			}
		}).Run(e.rconn)
	if err != nil {
		log.Error().Err(err).Msg("failed to query system stats table")
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

	var ts tableSize
	for cur.Next(&ts) {
		if cur.Err() != nil {
			log.Error().Err(cur.Err()).Msg("query error from cursor")
			errcount++
			return errcount
		}

		e.processTableSize(ts, ch)
	}

	return errcount
}

func (e *RethinkdbExporter) processTableSize(ts tableSize, ch chan<- prometheus.Metric) {
	ch <- prometheus.MustNewConstMetric(e.metrics.tableSize, prometheus.GaugeValue, ts.Size, ts.DB, ts.Table)
}
