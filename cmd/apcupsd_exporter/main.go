// Command apcupsd_exporter provides a Prometheus exporter for apcupsd.
package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/mdlayher/apcupsd"
	apcupsdexporter "github.com/mdlayher/apcupsd_exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	telemetryAddr = flag.String("telemetry.addr", ":9162", "address for apcupsd exporter")
	metricsPath   = flag.String("telemetry.path", "/metrics", "URL path for surfacing collected metrics")

	apcupsdAddr    = flag.String("apcupsd.addr", ":3551", "address of apcupsd Network Information Server (NIS)")
	apcupsdNetwork = flag.String("apcupsd.network", "tcp", `network of apcupsd Network Information Server (NIS): typically "tcp", "tcp4", or "tcp6"`)
)

func getenv_default(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func main() {
	usingEnvVars := os.Getenv("USE_ENV_VARS")
	if len(usingEnvVars) == 0 {
		log.Printf("USE_ENV_VARS not set, will parse flags")
		flag.Parse()
	} else {
		ev_telemetry_address := getenv_default("TELEMETRY_ADDRESS", ":9162")
		ev_metrics_path := getenv_default("METRICS_PATH", "/metrics")
		ev_apcupsdAddr := getenv_default("APCUPSD_ADDRESS", ":3551")
		ev_apcupsdNetwork := getenv_default("APCUPSD_NETWORK", "tcp")

		telemetryAddr = &ev_telemetry_address
		metricsPath = &ev_metrics_path
		apcupsdAddr = &ev_apcupsdAddr
		apcupsdNetwork = &ev_apcupsdNetwork

		log.Printf("Using Environment Variables: telemetryAddr: %s metricsPath: %s apcupsdAddr: %s apcupsdNetwork: %s",
			*telemetryAddr, *metricsPath, *apcupsdAddr, *apcupsdNetwork)
	}
	// flag.Parse()

	if *apcupsdAddr == "" {
		log.Fatal("address of apcupsd Network Information Server (NIS) must be specified with '-apcupsd.addr' flag")
	}

	fn := newClient(*apcupsdNetwork, *apcupsdAddr)

	prometheus.MustRegister(apcupsdexporter.New(fn))

	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsPath, http.StatusMovedPermanently)
	})

	log.Printf("starting apcupsd exporter on %q for server %s://%s",
		*telemetryAddr, *apcupsdNetwork, *apcupsdAddr)

	if err := http.ListenAndServe(*telemetryAddr, nil); err != nil {
		log.Fatalf("cannot start apcupsd exporter: %s", err)
	}
}

func newClient(network, addr string) apcupsdexporter.ClientFunc {
	return func(ctx context.Context) (*apcupsd.Client, error) {
		return apcupsd.DialContext(ctx, network, addr)
	}
}
