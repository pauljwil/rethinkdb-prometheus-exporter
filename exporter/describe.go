package exporter

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Describe sends metrics descriptions to the Prometheus channel.
func (e *RethinkdbExporter) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(e, ch)
}

// initMetrics initializes metrics descriptions.
func (e *RethinkdbExporter) initMetrics() {
	e.metrics.clusterClientConnections = prometheus.NewDesc(
		"cluster_client_connections",
		"Total number of connections from the cluster",
		nil, nil,
	)
	e.metrics.clusterDocsPerSecond = prometheus.NewDesc(
		"cluster_docs_per_second",
		"Total number of reads and writes of documents per second from the cluster",
		[]string{"operation"}, nil)

	e.metrics.serverClientConnections = prometheus.NewDesc(
		"server_client_connections",
		"Number of client connections to the server",
		[]string{"server"}, nil)
	e.metrics.serverQueriesPerSecond = prometheus.NewDesc(
		"server_queries_per_second",
		"Number of queries per second from the server",
		[]string{"server"}, nil)
	e.metrics.serverDocsPerSecond = prometheus.NewDesc(
		"server_docs_per_second",
		"Total number of reads and writes of documents per second from the server",
		[]string{"server", "operation"}, nil)

	e.metrics.tableDocsPerSecond = prometheus.NewDesc(
		"table_docs_per_second",
		"Number of reads and writes of documents per second from the table",
		[]string{"db", "table", "operation"}, nil)

	if e.collectTableStats {
		e.metrics.tableRowsCount = prometheus.NewDesc(
			"table_rows_count",
			"Approximate number of rows in the table",
			[]string{"db", "table"}, nil)
	}

	e.metrics.tableReplicaDocsPerSecond = prometheus.NewDesc(
		"tablereplica_docs_per_second",
		"Number of reads and writes of documents per second from the table replica",
		[]string{"db", "table", "server", "operation"}, nil)
	e.metrics.tableReplicaCacheBytes = prometheus.NewDesc(
		"tablereplica_cache_bytes",
		"Table replica cache size in bytes",
		[]string{"db", "table", "server"}, nil)
	e.metrics.tableReplicaIO = prometheus.NewDesc(
		"tablereplica_io",
		"Table replica reads and writes of bytes per second",
		[]string{"db", "table", "server", "operation"}, nil)
	e.metrics.tableReplicaDataBytes = prometheus.NewDesc(
		"tablereplica_data_bytes",
		"Table replica size in stored bytes",
		[]string{"db", "table", "server"}, nil)

	// Prometheus scrape metrics
	e.metrics.scrapeLatency = prometheus.NewDesc(
		"scrape_latency",
		"Latency of collecting scrape",
		nil, nil)
	e.metrics.scrapeErrors = prometheus.NewDesc(
		"scrape_errors",
		"Number of errors while collecting scrape",
		nil, nil)

	// current_issues table metrics
	e.metrics.logWriteIssues = prometheus.NewDesc(
		"log_write_issues",
		"Number of log write issues",
		nil, nil,
	)
	e.metrics.nameCollisionIssues = prometheus.NewDesc(
		"name_collision_issues",
		"Number of name collision issues",
		nil, nil,
	)
	e.metrics.outdatedIndexIssues = prometheus.NewDesc(
		"outdated_index_issues",
		"Number of outdated index issues",
		nil, nil,
	)
	e.metrics.totalAvailabilityIssues = prometheus.NewDesc(
		"total_availability_issues",
		"Number of total availability issues",
		nil, nil,
	)
	e.metrics.memoryAvailabilityIssues = prometheus.NewDesc(
		"memory_availability_issues",
		"Number of memory availability issues",
		nil, nil,
	)
	e.metrics.connectivityIssues = prometheus.NewDesc(
		"connectivity_issues",
		"Number of connectivity issues",
		nil, nil,
	)
	e.metrics.otherIssues = prometheus.NewDesc(
		"other_issues",
		"Number of unspecified issues",
		nil, nil,
	)

	// Table sizes
	e.metrics.tableSize = prometheus.NewDesc(
		"table_size",
		"RethinkDB table size in MB",
		[]string{"db", "table"}, nil,
	)
}
