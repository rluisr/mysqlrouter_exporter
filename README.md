mysqlrouter-exporter
=====================
[![Build Status](https://cloud.drone.io/api/badges/rluisr/mysqlrouter_exporter/status.svg)](https://cloud.drone.io/rluisr/mysqlrouter_exporter)

Usage
-----
1. Download binary from [release]().
2. Move to /usr/local/bin/
3. Add systemd script.
4. Start
```
[Unit]
Description=mysqlrouter-exporter
Documentation=https://github.com/rluisr/mysqlrouter-exporter
After=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/mysqlrouter_exporter --url mysqlrouter.luis.local --user luis --pass luis

[Install]
WantedBy=multi-user.target
```

You must set `--url`, `--user` and `--pass`.

Default listen port is `49152`.  
If want to change it, use `--port` flag.

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
WIP
