![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/elvin-tacirzade/log-exporter?logo=go)
[![Go Reference](https://pkg.go.dev/badge/github.com/elvin-tacirzade/log-exporter.svg)](https://pkg.go.dev/github.com/elvin-tacirzade/log-exporter)
![Docker Pulls](https://img.shields.io/docker/pulls/elvintacirzade/log-exporter?logo=docker&logoColor=white)
![Docker Image Size (tag)](https://img.shields.io/docker/image-size/elvintacirzade/log-exporter/latest?logo=docker&logoColor=white)

# Log Exporter

Log Exporter makes it possible to monitor the custom API logs using [Loki](https://grafana.com/oss/loki/).

## Overview

This project was inspired by [Promtail](https://grafana.com/docs/loki/latest/clients/promtail/). Your API have to run on [Docker](https://www.docker.com/). The exporter connects to your container, reads the container logs, pushes the logs to the [Loki](https://grafana.com/oss/loki/).

![Project schema](https://github.com/elvin-tacirzade/log-exporter/blob/main/grafana/assets/schema.png?raw=true)

## Getting Started

#### Log Structure
In the current version of the project, following a specific log structure is no longer mandatory. You can configure your logs with any fields you prefer.

However, there is an existing Grafana dashboard configured to work with the following log structure:

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
* `level` - log level.
* `method` - HTTP method.
* `status` - HTTP status code.
* `msg` - your log message.
* `dt` - device type.
* `timing` - your handle time (seconds) for each request.
* `ts` - log creation time.

The [Grafana](https://grafana.com/) dashboard includes specific visualizations for the `level` field, particularly for `error` and `warn` levels. If you use these log levels, corresponding graphs will be displayed on the dashboard, providing valuable insights into the occurrences and frequency of warnings and errors.

**Note**: The package you are using will take your log and add a `container` field before to push [Loki](https://grafana.com/oss/loki/).

### Usage Exporter

We use the following command to run it on [Docker](https://www.docker.com/).

```
docker run -d \
  --name log-exporter \
  --network monitoring \
  -v /var/run/docker.sock:/var/run/docker.sock \
  --env CONTAINER_NAME=your_container_name \
  --env LOKI_URL=http://loki:3100/loki/api/v1/push \
  --env DOCKER_RECONNECTION_TIME=20s \
  elvintacirzade/log-exporter:latest
```

[See](https://hub.docker.com/r/elvintacirzade/log-exporter) more information.

### Monitoring

We use [Loki](https://grafana.com/oss/loki/) and [Grafana](https://grafana.com/) for monitoring.

You must run [Loki](https://grafana.com/oss/loki/) and [Grafana](https://grafana.com/) on [Docker](https://www.docker.com/).


After running [Loki](https://grafana.com/oss/loki/) and [Grafana](https://grafana.com/) you must add [Loki](https://grafana.com/oss/loki/) data source in [Grafana](https://grafana.com/). Now you can import [dashboard](https://grafana.com/grafana/dashboards/19745-custom-logs/) for exporter.

![Grafana Dashboard Timing](https://github.com/elvin-tacirzade/log-exporter/blob/main/grafana/assets/dashboard-timing.png?raw=true)

![Grafana Dashboard Number of Processed Requests](https://github.com/elvin-tacirzade/log-exporter/blob/main/grafana/assets/dashboard-number-of-processed-requests.png?raw=true)

![Grafana Dashboard Info](https://github.com/elvin-tacirzade/log-exporter/blob/main/grafana/assets/dashboard-info.png?raw=true)
