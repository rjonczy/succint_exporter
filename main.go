package main

import (
	"net/http"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/prometheus/exporter-toolkit/web/kingpinflag"
)

var (
	metricsPath    = kingpin.Flag("metrics.path", "Path under which to expose Prometheus metrics.").Default("/metrics").String()
	succintProject = kingpin.Flag("succint.project", "Succint's project name").Default("@blasrodri/tendermintx-mainnet").String()
)

func init() {
	prometheus.MustRegister(version.NewCollector("succint_exporter"))
}

func main() {
	promlogConfig := &promlog.Config{}
	toolkitFlags := kingpinflag.AddFlags(kingpin.CommandLine, ":9103")
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("succint_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	logger := promlog.New(promlogConfig)

	level.Info(logger).Log("msg", "Starting succint_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "context", version.BuildContext())

	c := newSuccintCollector(logger)
	prometheus.MustRegister(c)

	http.Handle(*metricsPath, promhttp.Handler())
	if *metricsPath != "/" {

		landingConfig := web.LandingConfig{
			Name:        "succint_exporter",
			Description: "Prometheus Succint Exporter",
			Version:     version.Info(),
			Links: []web.LandingLinks{
				{
					Address: *metricsPath,
					Text:    "Metrics",
				},
			},
		}
		landingPage, err := web.NewLandingPage(landingConfig)
		if err != nil {
			level.Error(logger).Log("err", err)
			os.Exit(1)
		}
		http.Handle("/", landingPage)
	}

	srv := &http.Server{}
	if err := web.ListenAndServe(srv, toolkitFlags, logger); err != nil {
		level.Error(logger).Log("msg", "Error starting HTTP server", "err", err)
		os.Exit(1)
	}

}
