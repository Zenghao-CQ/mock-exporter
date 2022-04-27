package collector

import (
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// https://prometheus.io/docs/guides/go-application/#instrumenting-a-go-application-for-prometheus

// RecordMetrics add
func RecordMetrics(m *Metrics, start int64, period int64) {
	go func() {
		for {
			m.GenerateMockData()
			OpsProcessed.Inc()
			time.Sleep((time.Duration(start*1000 + rand.Int63n(period*1000))) * time.Microsecond)
		}
	}()
}

var (
	// OpsProcessed OpsProcessed
	OpsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name:        "failure_ops_total",
		Help:        "The total number of Mocking failures",
		ConstLabels: map[string]string{"app_name": "mock_exporter"},
		// ConstLabels: map[string]string{"test_node_exporter_label": "test_node_exporter_value", "app_name": "mock_exporter"},
	})
)

// func main() {
//         recordMetrics()

//         http.Handle("/metrics", promhttp.Handler())
//         http.ListenAndServe(":2112", nil)
// }
