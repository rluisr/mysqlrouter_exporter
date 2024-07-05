# mysqlrouter_exporter

[![lint](https://github.com/rluisr/mysqlrouter_exporter/actions/workflows/lint.yml/badge.svg)](https://github.com/rluisr/mysqlrouter_exporter/actions/workflows/lint.yml)
[![release](https://github.com/rluisr/mysqlrouter_exporter/actions/workflows/release.yml/badge.svg)](https://github.com/rluisr/mysqlrouter_exporter/actions/workflows/release.yml)

## Supported MySQL Router version

check [here](https://github.com/rluisr/mysqlrouter-go#supported-version)

## Usage

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

## Container

```bash
docker pull ghcr.io/rluisr/mysqlrouter_exporter:latest
```

[Packages](https://github.com/rluisr/mysqlrouter_exporter/pkgs/container/mysqlrouter_exporter)

## Environment

Edit systemd script or add an environment variables.

| Name                        | Default | Require | Description               |
| --------------------------- | ------- | ------- | ------------------------- |
| MYSQLROUTER_EXPORTER_URL    | -       | yes     | MySQL Router Rest API URL |
| MYSQLROUTER_EXPORTER_USER   | -       | no      | Username for REST API     |
| MYSQLROUTER_EXPORTER_PASS   | -       | no      | Password for REST API     |
| MYSQLROUTER_TLS_CACERT_PATH | -       | no      | TLS CA cert path          |
| MYSQLROUTER_TLS_CERT_PATH   | -       | no      | TLS cert path             |
| MYSQLROUTER_TLS_KEY_PATH    | -       | no      | TLS key path              |

You can also set it as a flag. See below.

```
Application Options:
      --url=                                                 MySQL Router Rest API URL [$MYSQLROUTER_EXPORTER_URL]
      --user=                                                Username for REST API [$MYSQLROUTER_EXPORTER_USER]
      --pass=                                                Password for REST API [$MYSQLROUTER_EXPORTER_PASS]
  -p, --listen-port=                                         Listen port (default: 9152)
      --service-name=                                        Service name for MySQL Router [$MYSQLROUTER_EXPORTER_SERVICE_NAME]
      --tls-ca-cert-path=                                    TLS CA cacert path [$MYSQLROUTER_TLS_CACERT_PATH]
      --tls-cert-path=                                       TLS cert path [$MYSQLROUTER_TLS_CERT_PATH]
      --tls-key-path=                                        TLS key path [$MYSQLROUTER_TLS_KEY_PATH]
  -k, --skip-tls-verify                                      Skip TLS Verification
      --collect-interval=                                    Collect interval time in sec. (default: 2)
      --collect.metadata.status                              Collect metrics from metadata status. CPU usage will increase.
      --collect.route.connections.byte_from_server           Collect metrics from route connections. CPU usage will increase.
      --collect.route.connections.byte_to_server             Collect metrics from route connections. CPU usage will increase.
      --collect.route.connections.time_started               Collect metrics from route connections. CPU usage will increase.
      --collect.route.connections.time_connected_to_server   Collect metrics from route connections. CPU usage will increase.
      --collect.route.connections.time_last_sent_to_server   Collect metrics from route connections. CPU usage will increase.
      --collect.route.connections.time_received_from_server  Collect metrics from route connections. CPU usage will increase.
  -v, --version                                              Show version

Help Options:
  -h, --help                                                 Show this help message
```

## Collector Flags

mysqlrouter_exporter can all get metrics. [MySQL Router REST API Reference](https://dev.mysql.com/doc/mysql-router/8.0/en/mysql-router-rest-api-reference.html)

| Name                                                | Default | Description                                                      |
| --------------------------------------------------- | ------- | ---------------------------------------------------------------- |
| collect.metadata.status                             | false   | Collect metrics from metadata status. CPU usage will increase.   |
| collect.route.connections.byte_from_server          | false   | Collect metrics from route connections. CPU usage will increase. |
| collect.route.connections.byte_to_server            | false   | Collect metrics from route connections. CPU usage will increase. |
| collect.route.connections.time_started              | false   | Collect metrics from route connections. CPU usage will increase. |
| collect.route.connections.time_connected_to_server  | false   | Collect metrics from route connections. CPU usage will increase. |
| collect.route.connections.time_last_sent_to_server  | false   | Collect metrics from route connections. CPU usage will increase. |
| collect.route.connections.time_received_from_server | false   | Collect metrics from route connections. CPU usage will increase. |

## Prometheus configuration

```yaml
scrape_configs:
  - job_name: "mysqlrouter"
    static_configs:
      - targets:
          - mysqlrouter.local:9152
```

## Grafana Dashboard

![Grafana Dashboard](img/grafana.png "Grafana Dashboard")

[Download dashboard](https://grafana.com/grafana/dashboards/10741)
