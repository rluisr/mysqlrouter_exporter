mysqlrouter_exporter
=====================
[![Build Status](https://cloud.drone.io/api/badges/rluisr/mysqlrouter_exporter/status.svg)](https://cloud.drone.io/rluisr/mysqlrouter_exporter)

Supported MySQL Router version
-------------------------------
8.0.16 - 8.0.20

Usage
-----
1. Enable REST API on your MySQL Router [here](https://lefred.be/content/mysqlrouter-8-0-17-and-the-rest-api/)
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
Environment="MYSQLROUTER_EXPORTER_URL=https://mysqlrouter-test.xzy.pw"
Environment="MYSQLROUTER_EXPORTER_USER=luis"
Environment="MYSQLROUTER_EXPORTER_PASS=luis"
ExecStart=/usr/local/bin/mysqlrouter_exporter

[Install]
WantedBy=multi-user.target
```

You must set these environment variables:  
- `MYSQLROUTER_EXPORTER_URL:` MySQL Router REST API URL.
- `MYSQLROUTER_EXPORTER_USER:` Username for REST API
- `MYSQLROUTER_EXPORTER_PASS:` Password for REST API


Default exporter listen port is `49152`. If you want change it, use flag `--port`.

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
