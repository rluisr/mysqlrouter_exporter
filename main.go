package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jessevdk/go-flags"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rluisr/mysqlrouter-go"
)

var (
	mysqlRouterClient *mysqlrouter.Client

	version string
	commit  string
	date    string
)

var args struct {
	RestAPIURL    string `short:"" long:"url" required:"true" env:"MYSQLROUTER_EXPORTER_URL" description:"MySQL Router Rest API URL"`
	RestAPIUser   string `short:"" long:"user" required:"false" env:"MYSQLROUTER_EXPORTER_USER" description:"Username for REST API"`
	RestAPIPass   string `short:"" long:"pass" required:"false" env:"MYSQLROUTER_EXPORTER_PASS" description:"Password for REST API"`
	ListenPort    int    `short:"p" long:"listen-port" default:"49152" description:"Listen port"`
	SkipTLSVerify bool   `short:"k" long:"skip-tls-verify" description:"Skip TLS Verification"`

	SkipCollectRouteConnectionsByteFromServer         bool `short:"" long:"skip.collect.route.connections.byte_from_server" description:"Skip Collect metrics from route connections. Set the flag if you getting high CPU usage."`
	SkipCollectRouteConnectionsByteToServer           bool `short:"" long:"skip.collect.route.connections.byte_to_server" description:"Skip Collect metrics from route connections. Set the flag if you getting high CPU usage."`
	SkipCollectRouteConnectionsTimeStarted            bool `short:"" long:"skip.collect.route.connections.time_started" description:"Skip Collect metrics from route connections. Set the flag if you getting high CPU usage."`
	SkipCollectRouteConnectionsTimeConnectedToServer  bool `short:"" long:"skip.collect.route.connections.time_connected_to_server" description:"Skip Collect metrics from route connections. Set the flag if you getting high CPU usage."`
	SkipCollectRouteConnectionsTimeLastSentToServer   bool `short:"" long:"skip.collect.route.connections.time_last_sent_to_server" description:"Skip Collect metrics from route connections. Set the flag if you getting high CPU usage."`
	SkipCollectRouteConnectionsTimeReceivedFromServer bool `short:"" long:"skip.collect.route.connections.time_received_from_server" description:"Skip Collect metrics from route connections. Set the flag if you getting high CPU usage."`

	Version bool `short:"v" long:"version" description:"Show version"`
}

const (
	nameSpace       = "mysqlrouter"
	collectInterval = 2 * time.Second
)

func initialClient() (*mysqlrouter.Client, error) {
	if args.RestAPIURL == "" {
		panic("These environments are missing.\n" +
			"MYSQLROUTER_EXPORTER_URL is required and MYSQLROUTER_EXPORTER_USER and MYSQLROUTER_EXPORTER_PASS are optional.")
	}

	return mysqlrouter.New(args.RestAPIURL, args.RestAPIUser, args.RestAPIPass, args.SkipTLSVerify)
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
			if !args.SkipCollectRouteConnectionsByteFromServer {
				routeConnectionsByteFromServerGauge.WithLabelValues(route.Name, router.Hostname, routeConnection.SourceAddress, routeConnection.DestinationAddress).Set(float64(routeConnection.BytesFromServer))
			}
			if !args.SkipCollectRouteConnectionsByteToServer {
				routeConnectionsByteToServerGauge.WithLabelValues(route.Name, router.Hostname, routeConnection.SourceAddress, routeConnection.DestinationAddress).Set(float64(routeConnection.BytesToServer))
			}
			if !args.SkipCollectRouteConnectionsTimeStarted {
				routeConnectionsTimeStartedGauge.WithLabelValues(route.Name, router.Hostname, routeConnection.SourceAddress, routeConnection.DestinationAddress).Set(float64(routeConnection.TimeStarted.Unix() * 1000)) // nolint
			}
			if !args.SkipCollectRouteConnectionsTimeConnectedToServer {
				routeConnectionsTimeConnectedToServerGauge.WithLabelValues(route.Name, router.Hostname, routeConnection.SourceAddress, routeConnection.DestinationAddress).Set(float64(routeConnection.TimeConnectedToServer.Unix() * 1000)) //nolint
			}
			if !args.SkipCollectRouteConnectionsTimeLastSentToServer {
				routeConnectionsTimeLastSentToServerGauge.WithLabelValues(route.Name, router.Hostname, routeConnection.SourceAddress, routeConnection.DestinationAddress).Set(float64(routeConnection.TimeLastSentToServer.Unix() * 1000)) // nolint
			}
			if !args.SkipCollectRouteConnectionsTimeReceivedFromServer {
				routeConnectionsTimeLastReceivedFromServerGauge.WithLabelValues(route.Name, router.Hostname, routeConnection.SourceAddress, routeConnection.DestinationAddress).Set(float64(routeConnection.TimeLastReceivedFromServer.Unix() * 1000)) // nolint
			}
		}
	}
}

func writeError(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "[mysqlrouter_exporter ERROR] %s\n", err.Error())
}

func main() {
	_, err := flags.Parse(&args)
	if err != nil {
		os.Exit(1)
	}
	if args.Version {
		fmt.Printf("version: %s commit: %s date: %s\n", version, commit, date)
		os.Exit(1)
	}

	mysqlRouterClient, err = initialClient()
	if err != nil {
		log.Fatalf("failed to create mysql router client. err: %s\n", err.Error())
	}

	recordMetrics()

	addr := fmt.Sprintf("0.0.0.0:%d", args.ListenPort)
	log.Printf("Start exporter on %s/metrics", addr)

	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe(addr, nil))
}
