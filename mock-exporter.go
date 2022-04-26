package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/zenghao-cq/mock-exporter/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// 命令行参数
	listenAddr       = flag.String("web.listen-port", "8000", "An port to listen on for web interface and telemetry.")
	metricsPath      = flag.String("web.telemetry-path", "/metrics", "A path under which to expose metrics.")
	metricsNamespace = flag.String("metric.namespace", "default", "Prometheus metrics namespace, as the prefix of metrics name")
	start            = flag.Int64("start", 60, "base of random time")
	period           = flag.Int64("period", 120, "period of random time")
	getPodPeriod     = flag.Int64("getPodPeriod", 10, "period(minutes) of refreash pod name")
	failureTypes     = flag.Int("failureTypes", 3, "types of failures")
)

func main() {

	flag.Parse()

	metrics := collector.NewMetrics(*metricsNamespace, *failureTypes)
	app := collector.NewApp(*metricsNamespace, *getPodPeriod)
	app.GetPodNames(metrics)
	collector.RecordMetrics(metrics, *start, *period)
	registry := prometheus.NewRegistry()
	registry.MustRegister(metrics, collector.OpsProcessed)

	http.Handle(*metricsPath, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
			<head><title>A Prometheus Exporter</title></head>
			<body>
			<h1>A Prometheus Exporter</h1>
			<p><a href='/metrics'>Metrics</a></p>
			</body>
			</html>`))
	})

	log.Printf("Starting Server at http://localhost:%s%s", *listenAddr, *metricsPath)
	log.Printf("start:%d period:%d failureTypes:%d", *start, *period, *failureTypes)
	log.Fatal(http.ListenAndServe(":"+*listenAddr, nil))
}
