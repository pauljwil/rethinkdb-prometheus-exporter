package exporter

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	gorethink "gopkg.in/rethinkdb/rethinkdb-go.v6"
)

var (
	cTables []issue
	sTables []stat

	sTable1 stat = stat{
		ID: []string{"cluster"},
		QueryEngine: queryEngine{
			ClientConnections: 43,
			ReadDocsPerSec:    0,
			WrittenDocsPerSec: 0,
		},
		Database: "dtr2",
		Table:    "cluster",
	}
	sTable2 stat = stat{
		ID: []string{"server"},
		QueryEngine: queryEngine{
			ClientConnections: 0,
			ReadDocsPerSec:    0,
			WrittenDocsPerSec: 0,
			QPS:               0,
		},
		Database: "dtr2",
		Table:    "server",
	}
	sTable3 stat = stat{
		ID: []string{"table"},
		QueryEngine: queryEngine{
			ReadDocsPerSec:    0,
			WrittenDocsPerSec: 0,
		},
		Database: "dtr2",
		Table:    "repository_team_access",
	}
	sTable4 stat = stat{
		ID: []string{"table_server"},
		QueryEngine: queryEngine{
			ReadDocsPerSec:    0,
			WrittenDocsPerSec: 0,
		},
		StorageEngine: storageEngine{
			Cache: cache{
				InUseBytes: 61384,
			},
			Disk: disk{
				SpaceUsage: spaceUsage{
					DataBytes: 2097152,
				},
			},
		},
		Database: "dtr2",
		Table:    "blob_links",
	}

	cTable1 issue = issue{
		Type: "log_write_error",
	}
	cTable2 issue = issue{
		Type: "server_name_collision",
	}
	cTable3 issue = issue{
		Type: "outdated_index",
	}
	cTable4 issue = issue{
		Type: "table_availability",
	}
	cTable5 issue = issue{
		Type: "memory_error",
	}
	cTable6 issue = issue{
		Type: "non_transitive_error",
	}
)

func TestNewStatsTableMetrics(t *testing.T) {
	sTables = append(sTables, sTable1, sTable2, sTable3, sTable4)
	cTables = append(cTables, cTable1, cTable2, cTable3, cTable4, cTable5, cTable6)

	mock := gorethink.NewMock()

	mock.On(gorethink.MockAnything()).Return(sTables, nil)
	mock.On(gorethink.MockAnything()).Return(cTables, nil)

	collector := &RethinkdbExporter{
		rconn:             mock,
		collectTableStats: true,
	}

	collector.initMetrics()

	prometheus.MustRegister(collector)

	server := httptest.NewServer(promhttp.Handler())
	defer server.Close()

	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("failed to get metrics: %s", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("unexpected status code: %d", resp.StatusCode)
	}

	expectedMetrics := []string{
		"cluster_client_connections",
		"cluster_docs_per_second",
		"server_client_connections",
		"server_queries_per_second",
		"server_docs_per_second",
		"table_docs_per_second",
		"table_rows_count",
		"tablereplica_docs_per_second",
		"tablereplica_cache_bytes",
		"tablereplica_io",
		"tablereplica_data_bytes",
		"scrape_latency",
		"scrape_errors",
		"log_write_issues",
		"name_collision_issues",
		"outdated_index_issues",
		"total_availability_issues",
		"memory_availability_issues",
		"connectivity_issues",
		"table_size",
	}

	bodyBytes, _ := io.ReadAll(resp.Body)

	bodyString := string(bodyBytes)

	for _, metricName := range expectedMetrics {
		if !strings.Contains(bodyString, metricName) {
			t.Errorf("expected metric %q not found in response", metricName)
		}
	}
}
