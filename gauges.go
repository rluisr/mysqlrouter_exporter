package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	routerStatusGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: nameSpace,
		Name:      "router_status",
		Help:      "MySQL Router information",
	}, []string{"process_id", "product_edition", "time_started", "version", "hostname"})
	metadataGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: nameSpace,
		Name:      "metadata",
		Help:      "metadata list",
	}, []string{"name"})
	metadataConfigGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: nameSpace,
		Name:      "metadata_config",
		Help:      "metadata config",
	}, []string{"name", "cluster_name", "time_refresh_in_ms", "group_replication_id"})
	metadataConfigNodeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: nameSpace,
		Name:      "metadata_config_node",
		Help:      "metadata config node",
	}, []string{"name", "router_host", "cluster_name", "hostname", "port"})
	metadataStatusGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: nameSpace,
		Name:      "metadata_status",
		Help:      "metadata status",
	}, []string{"name", "refresh_failed", "time_last_refresh_succeeded", "last_refresh_hostname", "last_refresh_port"})
	routeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: nameSpace,
		Name:      "route",
		Help:      "route name",
	}, []string{"name"})
	routeActiveConnectionsGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: nameSpace,
		Name:      "route_active_connections",
		Help:      "route active connections",
	}, []string{"name", "router_hostname"})
	routeTotalConnectionsGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: nameSpace,
		Name:      "route_total_connections",
		Help:      "route total connections",
	}, []string{"name", "router_hostname"})
	routeBlockedHostsGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: nameSpace,
		Name:      "route_blocked_hosts",
		Help:      "route blocked_hosts",
	}, []string{"name", "router_hostname"})
	routeHealthGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: nameSpace,
		Name:      "route_health",
		Help:      "0: not active, 1: active",
	}, []string{"name", "router_hostname"})
	routeDestinationsGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: nameSpace,
		Name:      "route_destinations",
		Help:      "",
	}, []string{"name", "address", "port"})
	routeConnectionsByteFromServerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "route_connections_byte_from_server",
		Help: "Route connections byte from server",
	}, []string{"name", "router_hostname", "source_address", "destination_address"})
	routeConnectionsByteToServerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "route_connections_byte_to_server",
		Help: "Route connections byte to server",
	}, []string{"name", "router_hostname", "source_address", "destination_address"})
	routeConnectionsTimeStartedGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "route_connections_time_started",
		Help: "Route connections time started",
	}, []string{"name", "router_hostname", "source_address", "destination_address"})
	routeConnectionsTimeConnectedToServerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "route_connections_time_connected_to_server",
		Help: "Route connections time connected to server",
	}, []string{"name", "router_hostname", "source_address", "destination_address"})
	routeConnectionsTimeLastSentToServerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "route_connections_time_last_sent_to_server",
		Help: "Route connections time last sent to server",
	}, []string{"name", "router_hostname", "source_address", "destination_address"})
	routeConnectionsTimeLastReceivedFromServerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "route_connections_time_last_received_from_server",
		Help: "Route connections time last received from server",
	}, []string{"name", "router_hostname", "source_address", "destination_address"})
)
