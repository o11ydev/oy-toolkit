# oy-expose

*oy-expose* reads a metrics file and exposes its content to be scraped by a
Prometheus server.

This is similar to the [Node Exporter Textfile Collector](https://github.com/prometheus/node_exporter#textfile-collector),
with a few differences:
- oy-expose only exposes on file.
- oy-expose does not embed other collectors.

## Example usage

Let's create a file called "metrics" with the following content:

```
maintenance_script_run_timestamp_seconds 1647524557
maintenance_script_return_code 0
```

We can run `oy-expose`:

```
$ oy-expose --web.disable-exporter-metrics
```

And query the metrics:

```
$ curl localhost:9099/metrics
maintenance_script_run_timestamp_seconds 1647524557
maintenance_script_return_code 0
```

