mysqlrouter_exporter
=====================
[![Build Status](https://cloud.drone.io/api/badges/rluisr/mysqlrouter_exporter/status.svg)](https://cloud.drone.io/rluisr/mysqlrouter_exporter)

Supported MySQL Router version
-------------------------------
check [here](https://github.com/rluisr/mysqlrouter-go#supported-version)

Usage
-----
1. Enable REST API on your MySQL Router [here](https://github.com/rluisr/mysqlrouter-go#supported-version)
2. Download binary from [release](https://github.com/rluisr/mysqlrouter_exporter/releases).
3. Move to /usr/local/bin/
4. Add systemd script.
5. Start
```
[Unit]
Description=mysqlrouter-exporter
Documentation=https://github.com/rluisr/mysqlrouter-exporter
After=network-online.target

[Service]
Type=simple
Environment="MYSQLROUTER_EXPORTER_URL=url"
Environment="MYSQLROUTER_EXPORTER_USER=user"
Environment="MYSQLROUTER_EXPORTER_PASS=pass"
ExecStart=/usr/local/bin/mysqlrouter_exporter

[Install]
WantedBy=multi-user.target
```

Environment
-----------

Edit systemd script or add an environment variables.

Name                        | Default | Require | Description
----------------------------|---------| --------| ----------
MYSQLROUTER_EXPORTER_URL    | -       | yes     | MySQL Router Rest API URL
MYSQLROUTER_EXPORTER_USER   | -       | no      | Username for REST API
MYSQLROUTER_EXPORTER_PASS   | -       | no      | Password for REST API

You can also set it as a flag.

Collector Flags
----------------

Name                                                        | Default   | Description
------------------------------------------------------------|-----------|-------------
skip.collect.route.connections.byte_from_server             | false     | Skip Collect metrics from route connections. Set the flag if you getting high CPU usage.
skip.collect.route.connections.byte_to_server               | false     | Skip Collect metrics from route connections. Set the flag if you getting high CPU usage.
skip.collect.route.connections.time_started                 | false     | Skip Collect metrics from route connections. Set the flag if you getting high CPU usage.
skip.collect.route.connections.time_connected_to_server     | false     | Skip Collect metrics from route connections. Set the flag if you getting high CPU usage.
skip.collect.route.connections.time_last_sent_to_server     | false     | Skip Collect metrics from route connections. Set the flag if you getting high CPU usage.
skip.collect.route.connections.time_received_from_server    | false     | Skip Collect metrics from route connections. Set the flag if you getting high CPU usage.

```
Usage:
  mysqlrouter_exporter [OPTIONS]

Application Options:
      --url=                                                      MySQL Router Rest API URL [$MYSQLROUTER_EXPORTER_URL]
      --user=                                                     Username for REST API [$MYSQLROUTER_EXPORTER_USER]
      --pass=                                                     Password for REST API [$MYSQLROUTER_EXPORTER_PASS]
  -p, --listen-port=                                              Listen port (default: 49152)
  -k, --skip-tls-verify                                           Skip TLS Verification
      --skip.collect.route.connections.byte_from_server           Skip Collect metrics from route connections. Set the flag if you getting high CPU
                                                                  usage.
      --skip.collect.route.connections.byte_to_server             Skip Collect metrics from route connections. Set the flag if you getting high CPU
                                                                  usage.
      --skip.collect.route.connections.time_started               Skip Collect metrics from route connections. Set the flag if you getting high CPU
                                                                  usage.
      --skip.collect.route.connections.time_connected_to_server   Skip Collect metrics from route connections. Set the flag if you getting high CPU
                                                                  usage.
      --skip.collect.route.connections.time_last_sent_to_server   Skip Collect metrics from route connections. Set the flag if you getting high CPU
                                                                  usage.
      --skip.collect.route.connections.time_received_from_server  Skip Collect metrics from route connections. Set the flag if you getting high CPU
                                                                  usage.
  -v, --version                                                   Show version

Help Options:
  -h, --help                                                      Show this help message
```

Prometheus configuration
-------------------------
```yaml
scrape_configs:
  - job_name: 'mysqlrouter'
    static_configs:
      - targets:
        - mysqlrouter01.luis.local:49152
```

Grafana Dashboard
------------------------
![Grafana Dashboard](https://grafana.com/api/dashboards/10741/images/6783/image "Grafana Dashboard")

available [here](https://grafana.com/grafana/dashboards/10741).
