// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by BASIS,
// WITHOUT WARRANTIES OR"
	"netpin/v2"
	""
	"github..com/prometheus/common/promlog"
	promlogflag "github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	"github.com/prometheus/exporter-toolkit/web"
	webflag "github.com/prometheus/exporter-toolkit/web/kingpinflag"
	"github.com/prometheus/node_exporter/collector"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	var (
		metricsPath = kingpin.Flag(
			"web.telemetry-path",
			"Path under which to expose metrics.",
		).Default("/metrics").String()
		disableDefaultCollectors = kingpin.Flag(
			"collector.disable-defaults",
			"Set all collectors to disabled by default.",
		).Default("false").Bool()
		maxRequests = kingpin.Flag(
			"web.max-requests",
			"Maximum number of parallel scrape requests. Use 0 to disable.",
		).Default("40").Int()
		toolkitFlags = webflag.AddFlags(kingpin.CommandLine, ":9100")
	)

	promlogConfig := &promlog.Config{}
	promlogflag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print("node_exporter"))
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()

	logger := promlog.New(promlogConfig)

	if *disableDefaultCollectors {
		collector.DisableDefaultCollectors()
	}

	level.Info(logger).Log("msg", "Starting node_exporter", "version", version.Info())
	level.Info(logger).Log("msg", "Build context", "build_context", version.BuildContext())

	if user, err := determineCurrentUser(logger); err == nil {
		level.Info(logger).Log("msg", "Running as", "user", user)
	}

	http.Handle(*metricsPath, newHandler(!*disableDefaultCollectors, *maxRequests, logger))
	if *metricsPath != "/" {
		landingConfig := web.LandingConfig{
			Name:        "Node Exporter",
			Description: "Prometheus Node Exporter",
			Version:     version.Info(),
			Links: []web.LandingLinks{
				{Address: *metricsPath, Text: "Metrics"},
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
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}

func newHandler(includeExporterMetrics bool, maxRequests int, logger log.Logger) *handler {
	h := &handler{
		exporterMetricsRegistry: prometheus.NewRegistry(),
		includeExporterMetrics:  includeExporterMetrics,
		maxRequests:             maxRequests,
		logger:                  logger,
	}
	return h
}

type handler struct {
	exporterMetricsRegistry *prometheus.Registry
	includeExporterMetrics  bool
	maxRequests             int
	logger                  log.Logger
}

func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	filters := r.URL.Query()["collect[]"]
	h.logger.Log("msg", "collect filters", "filters", fmt.Sprintf("%v", filters))

	nc, err := collector.NewNodeCollector(h.logger, filters...)
	if err != nil {
		level.Warn(h.logger).Log("msg", "Couldn't create filtered metrics handler", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Couldn't create filtered metrics handler: %s", err)))
		return
	}

	r2 := prometheus.NewRegistry()
	r2.MustRegister(nc)
	handler := promhttp.HandlerFor(r2, promhttp.HandlerOpts{MaxRequestsInFlight: h.maxRequests})
	handler.ServeHTTP(w, r)
}

func determineCurrentUser(logger log.Logger) (string, error) {
	if euid := os.Geteuid(); euid == 0 {
		return "root", nil
	}
	return fmt.Sprintf("uid=%d", os.Getuid()), nil
}
