package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rluisr/mysqlrouter-go"
)

var nameSpace = "mysqlrouter"

var (
	routerStatusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "router_status",
			Help:      "MySQL Router information",
		}, []string{"process_id", "product_edition", "time_started", "version", "hostname"})
	metadataGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "metadata",
			Help:      "metadata list",
		}, []string{"name"})
	metadataConfigGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "metadata_config",
			Help:      "metadata config",
		}, []string{"name", "cluster_name", "time_refresh_in_ms", "group_replication_id"})
	metadataConfigNodeGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "metadata_config_node",
			Help:      "metadata config node",
		}, []string{"name", "cluster_name", "hostname", "port"})
	metadataStatusGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "metadata_status",
			Help:      "metadata status",
		}, []string{"name", "refresh_failed", "time_last_refresh_succeeded", "last_refresh_hostname", "last_refresh_port"})
	routeGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route",
			Help:      "route name",
		}, []string{"name"})
	routeActiveConnectionsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_active_connections",
			Help:      "route active connections",
		}, []string{"name"})
	routeTotalConnectionsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_total_connections",
			Help:      "route total connections",
		}, []string{"name"})
	routeBlockedHostsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_blocked_hosts",
			Help:      "route blocked_hosts",
		}, []string{"name"})
	routeHealthGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_health",
			Help:      "0: not active, 1: active",
		}, []string{"name"})
	routeDestinationsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: nameSpace,
			Name:      "route_destinations",
			Help:      "",
		}, []string{"name", "address", "port"})
	/* WIP
	routeConnectionsGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "route_connections",
			Help: "",
		}, []string{"name", "", "port"})
	*/
)

func init() {
	prometheus.MustRegister(
		routerStatusGauge,
		metadataGauge,
		metadataConfigGauge,
		metadataConfigNodeGauge,
		metadataStatusGauge,
		routeGauge,
		routeActiveConnectionsGauge,
		routeTotalConnectionsGauge,
		routeBlockedHostsGauge,
		routeHealthGauge,
		routeDestinationsGauge,
	)
}

var (
	port = flag.String("port", "49152", "Listen port for exporter")
	url  = flag.String("url", "", "URL of MySQL Router REST API")
	user = flag.String("user", "", "Username for MySQL Router REST API")
	pass = flag.String("pass", "", "Password for MySQL Router REST API")
)

func main() {
	flag.Parse()

	if *url == "" || *user == "" || *pass == "" {
		panic("--url, --user and --pass is must be set.")
	}

	mr, err := mysqlrouter.New("https://mysqlrouter-test.xzy.pw", "luis", "luis")
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			// Router
			router, err := mr.GetRouterStatus()
			if err != nil {
				panic(err)
			}
			routerStatusGauge.WithLabelValues(strconv.Itoa(router.ProcessID), router.ProductEdition, router.TimeStarted.String(), router.Version, router.Hostname)

			// Metadata
			metadata, err := mr.GetAllMetadata()
			if err != nil {
				panic(err)
			}
			for _, m := range metadata {
				metadataGauge.WithLabelValues(m.Name)

				// config
				mc, err := mr.GetMetadataConfig(m.Name)
				if err != nil {
					panic(err)
				}
				metadataConfigGauge.WithLabelValues(m.Name, mc.ClusterName, strconv.Itoa(mc.TimeRefreshInMs), mc.GroupReplicationID)

				// config node
				for _, node := range mc.Nodes {
					metadataConfigNodeGauge.WithLabelValues(m.Name, mc.ClusterName, node.Hostname, strconv.Itoa(node.Port))
				}

				// status
				ms, err := mr.GetMetadataStatus(m.Name)
				if err != nil {
					panic(err)
				}
				metadataStatusGauge.WithLabelValues(m.Name, strconv.Itoa(ms.RefreshFailed), ms.TimeLastRefreshSucceeded.String(), ms.LastRefreshHostname, strconv.Itoa(ms.LastRefreshPort))
			}

			// Routes
			routes, err := mr.GetAllRoutes()
			if err != nil {
				panic(err)
			}
			for _, route := range routes {
				routeGauge.WithLabelValues(route.Name)

				rs, err := mr.GetRouteStatus(route.Name)
				if err != nil {
					panic(err)
				}
				routeActiveConnectionsGauge.WithLabelValues(route.Name).Set(float64(rs.ActiveConnections))
				routeTotalConnectionsGauge.WithLabelValues(route.Name).Set(float64(rs.TotalConnections))
				routeBlockedHostsGauge.WithLabelValues(route.Name).Set(float64(rs.BlockedHosts))

				rh, err := mr.GetRouteHealth(route.Name)
				if err != nil {
					panic(err)
				}
				if rh.IsAlive {
					routeHealthGauge.WithLabelValues(route.Name).Set(float64(1))
				} else {
					routeHealthGauge.WithLabelValues(route.Name).Set(float64(0))
				}

				rd, err := mr.GetRouteDestinations(route.Name)
				if err != nil {
					panic(err)
				}
				for _, d := range rd {
					routeDestinationsGauge.WithLabelValues(route.Name, d.Address, strconv.Itoa(d.Port))
				}

				/* WIP
				rc, err := mr.GetRouteConnections(route.Name)
				if err != nil {
					panic(err)
				}
				*/
			}
			time.Sleep(60 * time.Second)
		}
	}()

	log.Printf("listen: %s\n", "0.0.0.0:"+*port)
	http.Handle("/metrics", promhttp.Handler())
	log.Fatal(http.ListenAndServe("0.0.0.0:"+*port, nil))
}
