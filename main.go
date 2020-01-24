package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rluisr/mysqlrouter-go"
)

var (
	nameSpace = "mysqlrouter"
	metricMap map[string]*prometheus.GaugeVec
	port      = os.Getenv("MYSQLROUTER_EXPORTER_PORT")
	url       = os.Getenv("MYSQLROUTER_EXPORTER_URL")
	user      = os.Getenv("MYSQLROUTER_EXPORTER_USER")
	pass      = os.Getenv("MYSQLROUTER_EXPORTER_PASS")
)

func init() {
	metricMap = map[string]*prometheus.GaugeVec{
		"router status": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "router_status",
			Help:      "MySQL Router information",
		}, []string{"process_id", "product_edition", "time_started", "version", "hostname"}),
		"metadata": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "metadata",
			Help:      "metadata list",
		}, []string{"name"}),
		"metadata config": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "metadata_config",
			Help:      "metadata config",
		}, []string{"name", "cluster_name", "time_refresh_in_ms", "group_replication_id"}),
		"metadata config node": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "metadata_config_node",
			Help:      "metadata config node",
		}, []string{"name", "router_host", "cluster_name", "hostname", "port"}),
		"metadata status": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "metadata_status",
			Help:      "metadata status",
		}, []string{"name", "refresh_failed", "time_last_refresh_succeeded", "last_refresh_hostname", "last_refresh_port"}),
		"route": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route",
			Help:      "route name",
		}, []string{"name"}),
		"route active connections": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_active_connections",
			Help:      "route active connections",
		}, []string{"name", "router_hostname"}),
		"route total connections": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_total_connections",
			Help:      "route total connections",
		}, []string{"name", "router_hostname"}),
		"route blocked hosts": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_blocked_hosts",
			Help:      "route blocked_hosts",
		}, []string{"name", "router_hostname"}),
		"route health": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_health",
			Help:      "0: not active, 1: active",
		}, []string{"name", "router_hostname"}),
		"route destinations": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_destinations",
			Help:      "",
		}, []string{"name", "address", "port"}),
		"route connections bytes from server": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_connections_byte_from_server",
			Help:      "Route connections byte from server",
		}, []string{"name", "router_hostname", "source_address", "destination_address"}),
		"route connections byte to server": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_connections_byte_to_server",
			Help:      "Route connections byte to server",
		}, []string{"name", "router_hostname", "source_address", "destination_address"}),
		"route connections time started": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_connections_time_started",
			Help:      "Route connections time started",
		}, []string{"name", "router_hostname", "source_address", "destination_address"}),
		"route connections time connected to server": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_connections_time_connected_to_server",
			Help:      "Route connections time connected to server",
		}, []string{"name", "router_hostname", "source_address", "destination_address"}),
		"route connections time last sent to server": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_connections_time_last_sent_to_server",
			Help:      "Route connections time last sent to server",
		}, []string{"name", "router_hostname", "source_address", "destination_address"}),
		"route connections time last received from server": prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_connections_time_last_received_from_server",
			Help:      "Route connections time last received from server",
		}, []string{"name", "router_hostname", "source_address", "destination_address"}),
	}
	for _, metric := range metricMap {
		prometheus.MustRegister(metric)
	}
}

func main() {
	flag.Parse()

	if url == "" || user == "" || pass == "" {
		panic("The environment missing.\n" +
			"MYSQLROUTER_EXPORTER_URL, MYSQLROUTER_EXPORTER_USER and MYSQLROUTER_EXPORTER_PASS is required.")
	}

	mr, err := mysqlrouter.New(url, user, pass)
	if err != nil {
		writeError(err)
	}

	go func() {
		for {
			// Router
			router, err := mr.GetRouterStatus()
			if err != nil {
				writeError(err)
			}
			metricMap["router status"].WithLabelValues(strconv.Itoa(router.ProcessID), router.ProductEdition, router.TimeStarted.String(), router.Version, router.Hostname)

			// Metadata
			metadata, err := mr.GetAllMetadata()
			if err != nil {
				writeError(err)
			}
			for _, m := range metadata {
				metricMap["metadata"].WithLabelValues(m.Name)

				// config
				mc, err := mr.GetMetadataConfig(m.Name)
				if err != nil {
					writeError(err)
				}
				metricMap["metadata config"].WithLabelValues(m.Name, mc.ClusterName, strconv.Itoa(mc.TimeRefreshInMs), mc.GroupReplicationID)

				// config node
				for _, node := range mc.Nodes {
					metricMap["metadata config node"].WithLabelValues(m.Name, router.Hostname, mc.ClusterName, node.Hostname, strconv.Itoa(node.Port))
				}

				// status
				ms, err := mr.GetMetadataStatus(m.Name)
				if err != nil {
					writeError(err)
				}
				metricMap["metadata status"].WithLabelValues(m.Name, strconv.Itoa(ms.RefreshFailed), ms.TimeLastRefreshSucceeded.String(), ms.LastRefreshHostname, strconv.Itoa(ms.LastRefreshPort))
			}

			// Routes
			routes, err := mr.GetAllRoutes()
			if err != nil {
				writeError(err)
			}
			for _, route := range routes {
				metricMap["route"].WithLabelValues(route.Name)

				rs, err := mr.GetRouteStatus(route.Name)
				if err != nil {
					writeError(err)
				}
				metricMap["route active connections"].WithLabelValues(route.Name, router.Hostname).Set(float64(rs.ActiveConnections))
				metricMap["route total connections"].WithLabelValues(route.Name, router.Hostname).Set(float64(rs.TotalConnections))
				metricMap["route blocked hosts"].WithLabelValues(route.Name, router.Hostname).Set(float64(rs.BlockedHosts))

				rh, err := mr.GetRouteHealth(route.Name)
				if err != nil {
					writeError(err)
				}
				if rh.IsAlive {
					metricMap["route health"].WithLabelValues(route.Name, router.Hostname).Set(float64(1))
				} else {
					metricMap["route health"].WithLabelValues(route.Name).Set(float64(0))
				}

				rd, err := mr.GetRouteDestinations(route.Name)
				if err != nil {
					writeError(err)
				}
				for _, d := range rd {
					metricMap["route destinations"].WithLabelValues(route.Name, d.Address, strconv.Itoa(d.Port))
				}

				rc, err := mr.GetRouteConnections(route.Name)
				if err != nil {
					writeError(err)
				}
				for _, c := range rc {
					metricMap["route connections bytes from server"].WithLabelValues(route.Name, router.Hostname, c.SourceAddress, c.DestinationAddress).Set(float64(c.BytesFromServer))
					metricMap["route connections byte to server"].WithLabelValues(route.Name, router.Hostname, c.SourceAddress, c.DestinationAddress).Set(float64(c.BytesToServer))
					metricMap["route connections time started"].WithLabelValues(route.Name, router.Hostname, c.SourceAddress, c.DestinationAddress).Set(float64(c.TimeStarted.Unix() * 1000))
					metricMap["route connections time connected to server"].WithLabelValues(route.Name, router.Hostname, c.SourceAddress, c.DestinationAddress).Set(float64(c.TimeConnectedToServer.Unix() * 1000))
					metricMap["route connections time last sent to server"].WithLabelValues(route.Name, router.Hostname, c.SourceAddress, c.DestinationAddress).Set(float64(c.TimeLastSentToServer.Unix() * 1000))
					metricMap["route connections time last received from server"].WithLabelValues(route.Name, router.Hostname, c.SourceAddress, c.DestinationAddress).Set(float64(c.TimeLastReceivedFromServer.Unix() * 1000))
				}
			}
			time.Sleep(60 * time.Second)
		}
	}()

	if port == "" {
		port = "49152"
	}

	log.Printf("listen: %s\n", "0.0.0.0:"+port)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

func writeError(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "[mysqlrouter_exporter ERROR] %v\n", err)
}
