![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/elvin-tacirzade/log-exporter?logo=go)
[![Go Reference](https://pkg.go.dev/badge/github.com/elvin-tacirzade/log-exporter.svg)](https://pkg.go.dev/github.com/elvin-tacirzade/log-exporter)
![Docker Pulls](https://img.shields.io/docker/pulls/elvintacirzade/log-exporter?logo=docker&logoColor=white)
![Docker Image Size (tag)](https://img.shields.io/docker/image-size/elvintacirzade/log-exporter/latest?logo=docker&logoColor=white)

# Log Exporter

Log Exporter makes it possible to monitor the custom API logs using [Loki](https://grafana.com/oss/loki/).

## Overview

This project was inspired by [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/). Your API have to run on [Docker](https://www.docker.com/). The exporter connects to your container, reads the container logs, pushes the logs to the [Loki](https://grafana.com/oss/loki/).

![Project schema](https://github.com/elvin-tacirzade/log-exporter/blob/main/grafana/photos/schema.png?raw=true)

## Getting Started

We need some prerequisites to start exporting.

### Prerequisites

#### Log Structure
First of all your api log must be like this:

```
{"ip":"192.168.1.1","caller":"app/main.go:102","path":"/users","level":"info","method":"GET","status":200,"msg":"get users successfully","dt":"mobile","timing":0.776347977,"ts":"2023-07-10T13:01:38Z"}
```

```
{
    "ip":"192.168.1.1",
    "caller":"app/main.go:102",
    "path":"/users",
    "level":"info",
    "method":"GET",
    "status":200,
    "msg":"get users successfully",
    "dt":"mobile",
    "timing":0.776347977,
    "ts":"2023-07-10T13:01:38Z"
}
```

* `ip` - ip address.
* `caller` - log line in your application.
* `path` - your api path.
* `level` - log level. It can be `info`, `error`, `warn`, `fatal`, `debug` and `trace`.
* `method` - HTTP method.
* `status` - HTTP status code and must be `integer`.
* `msg` - your log message.
* `dt` - device type.
* `timing` - your handle time for each request. It must be `float`.
* `ts` - log creation time.

### Usage Exporter

We use the following command to run it on [Docker](https://www.docker.com/).

```
docker run -d \
  --name log-exporter \
  --network monitoring \
  -v /var/lib/docker/containers/8528e90660bce465a47d381d4dd8d5d43ef5380d0d35a0ea28d1fb85382561a3:/logs \
  --env CONTAINER_ID=8528e90660bce465a47d381d4dd8d5d43ef5380d0d35a0ea28d1fb85382561a3 \
  --env LOKI_URL=http://loki:3100/loki/api/v1/push \
  elvintacirzade/log-exporter:latest
```

**Don't forget to change volume mapping and environment variables.**

[See](https://hub.docker.com/r/elvintacirzade/log-exporter) more information.

### Monitoring

We use [Loki](https://grafana.com/oss/loki/) and [Grafana](https://grafana.com/) for monitoring.

You must run [Loki](https://grafana.com/oss/loki/) and [Grafana](https://grafana.com/) on [Docker](https://www.docker.com/).


After running [Loki](https://grafana.com/oss/loki/) and [Grafana](https://grafana.com/) you must add [Loki](https://grafana.com/oss/loki/) data source in [Grafana](https://grafana.com/). Now you can import [dashboard](https://grafana.com/grafana/dashboards/19199-logs/) for exporter.

![Grafana Dashboard](https://github.com/elvin-tacirzade/log-exporter/blob/main/grafana/photos/dashboard_1.png?raw=true)

![Grafana Dashboard](https://github.com/elvin-tacirzade/log-exporter/blob/main/grafana/photos/dashboard_2.png?raw=true)

![Grafana Dashboard](https://github.com/elvin-tacirzade/log-exporter/blob/main/grafana/photos/dashboard_3.png?raw=true)
