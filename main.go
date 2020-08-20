package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rluisr/mysqlrouter-go"
)

var (
	mysqlRouterClient *mysqlrouter.Client
	url               = os.Getenv("MYSQLROUTER_EXPORTER_URL")
	user              = os.Getenv("MYSQLROUTER_EXPORTER_USER")
	pass              = os.Getenv("MYSQLROUTER_EXPORTER_PASS")

	version string
	commit  string
	date    string
)

const (
	nameSpace       = "mysqlrouter"
	collectInterval = 2 * time.Second
)

func initialClient(skipTLSVerify bool) (*mysqlrouter.Client, error) {
	if url == "" {
		panic("These environments are missing.\n" +
			"MYSQLROUTER_EXPORTER_URL is required and MYSQLROUTER_EXPORTER_USER and MYSQLROUTER_EXPORTER_PASS are optional.")
	}

	return mysqlrouter.New(url, user, pass, skipTLSVerify)
}

func recordMetrics() {
	go func() {
		for {
			collectMetrics()
			time.Sleep(collectInterval)
		}
	}()
}

func collectMetrics() {
	// router
	router, err := mysqlRouterClient.GetRouterStatus()
	if err != nil {
		writeError(err)
		return
	}
	routerStatusGauge.WithLabelValues(strconv.Itoa(router.ProcessID), router.ProductEdition, router.TimeStarted.String(), router.Version, router.Hostname)

	// metadata
	metadatas, err := mysqlRouterClient.GetAllMetadata()
	if err != nil {
		writeError(err)
		return
	}
	for _, metadata := range metadatas {
		metadataGauge.WithLabelValues(metadata.Name)

		// config
		metadataConfig, gmcErr := mysqlRouterClient.GetMetadataConfig(metadata.Name)
		if gmcErr != nil {
			writeError(gmcErr)
			return
		}
		metadataConfigGauge.WithLabelValues(metadata.Name, metadataConfig.ClusterName, strconv.Itoa(metadataConfig.TimeRefreshInMs), metadataConfig.GroupReplicationID)

		// config node
		for _, metadataConfigNode := range metadataConfig.Nodes {
			metadataConfigNodeGauge.WithLabelValues(metadata.Name, router.Hostname, metadataConfig.ClusterName, metadataConfigNode.Hostname, strconv.Itoa(metadataConfigNode.Port))
		}

		// status
		metadataStatus, gmsErr := mysqlRouterClient.GetMetadataStatus(metadata.Name)
		if gmsErr != nil {
			writeError(gmsErr)
			return
		}
		metadataStatusGauge.WithLabelValues(metadata.Name, strconv.Itoa(metadataStatus.RefreshFailed), metadataStatus.TimeLastRefreshSucceeded.String(), metadataStatus.LastRefreshHostname, strconv.Itoa(metadataStatus.LastRefreshPort))
	}

	// routes
	routes, err := mysqlRouterClient.GetAllRoutes()
	if err != nil {
		writeError(err)
		return
	}
	for _, route := range routes {
		routeGauge.WithLabelValues(route.Name)

		routeStatus, err := mysqlRouterClient.GetRouteStatus(route.Name)
		if err != nil {
			writeError(err)
			return
		}

		routeActiveConnectionsGauge.WithLabelValues(route.Name, router.Hostname).Set(float64(routeStatus.ActiveConnections))
		routeTotalConnectionsGauge.WithLabelValues(route.Name, router.Hostname).Set(float64(routeStatus.TotalConnections))
		routeBlockedHostsGauge.WithLabelValues(route.Name, router.Hostname).Set(float64(routeStatus.BlockedHosts))

		routeHealth, err := mysqlRouterClient.GetRouteHealth(route.Name)
		if err != nil {
			writeError(err)
			return
		}
		if routeHealth.IsAlive {
			routeHealthGauge.WithLabelValues(route.Name, router.Hostname).Set(float64(1))
		} else {
			routeHealthGauge.WithLabelValues(route.Name).Set(float64(0))
		}

		routeDestinations, err := mysqlRouterClient.GetRouteDestinations(route.Name)
		if err != nil {
			writeError(err)
			return
		}
		for _, routeDestination := range routeDestinations {
			routeDestinationsGauge.WithLabelValues(route.Name, routeDestination.Address, strconv.Itoa(routeDestination.Port))
		}

		routeConnections, err := mysqlRouterClient.GetRouteConnections(route.Name)
		if err != nil {
			writeError(err)
			return
		}
		for _, routeConnection := range routeConnections {
			routeConnectionsByteFromServerGauge.WithLabelValues(route.Name, router.Hostname, routeConnection.SourceAddress, routeConnection.DestinationAddress).Set(float64(routeConnection.BytesFromServer))
			routeConnectionsByteToServerGauge.WithLabelValues(route.Name, router.Hostname, routeConnection.SourceAddress, routeConnection.DestinationAddress).Set(float64(routeConnection.BytesToServer))
			routeConnectionsTimeStartedGauge.WithLabelValues(route.Name, router.Hostname, routeConnection.SourceAddress, routeConnection.DestinationAddress).Set(float64(routeConnection.TimeStarted.Unix() * 1000))                               // nolint
			routeConnectionsTimeConnectedToServerGauge.WithLabelValues(route.Name, router.Hostname, routeConnection.SourceAddress, routeConnection.DestinationAddress).Set(float64(routeConnection.TimeConnectedToServer.Unix() * 1000))           //nolint
			routeConnectionsTimeLastSentToServerGauge.WithLabelValues(route.Name, router.Hostname, routeConnection.SourceAddress, routeConnection.DestinationAddress).Set(float64(routeConnection.TimeLastSentToServer.Unix() * 1000))             // nolint
			routeConnectionsTimeLastReceivedFromServerGauge.WithLabelValues(route.Name, router.Hostname, routeConnection.SourceAddress, routeConnection.DestinationAddress).Set(float64(routeConnection.TimeLastReceivedFromServer.Unix() * 1000)) // nolint
		}
	}
}

func writeError(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "[mysqlrouter_exporter ERROR] %s\n", err.Error())
}

func flagUsage() {
	usageText := `--port        		Listen port. Default 49152
--version     		Show version
--skip-tls-verify	Skip TLS Verification`

	_, _ = fmt.Fprintf(os.Stderr, "%s\n\n", usageText)
}

func main() {
	listenPortFlag := flag.Int("port", 49152, "listen port")
	versionFlag := flag.Bool("version", false, "show version --version")
	skipTLSVerify := flag.Bool("skip-tls-verify", false, "Skip TLS Verification")

	flag.Usage = flagUsage
	flag.Parse()

	if *versionFlag {
		fmt.Printf("version: %s commit: %s date: %s\n", version, commit, date)
		os.Exit(0)
	}

	var err error
	mysqlRouterClient, err = initialClient(*skipTLSVerify)
	if err != nil {
		log.Fatalf("failed to create mysql router client. err: %s\n", err.Error())
	}

	recordMetrics()

	addr := fmt.Sprintf("0.0.0.0:%d", *listenPortFlag)
	log.Printf("Start exporter on %s/metrics", addr)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}
