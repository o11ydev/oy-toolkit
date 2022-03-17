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
# HELP maintenance_script_return_code Metric read from metrics
# TYPE maintenance_script_return_code untyped
maintenance_script_return_code 0
# HELP maintenance_script_run_timestamp_seconds Metric read from metrics
# TYPE maintenance_script_run_timestamp_seconds untyped
maintenance_script_run_timestamp_seconds 1.647524557e+09
# HELP node_textfile_mtime_seconds Unixtime mtime of textfiles successfully read.
# TYPE node_textfile_mtime_seconds gauge
node_textfile_mtime_seconds{file="metrics"} 1.647530635e+09
# HELP node_textfile_scrape_error 1 if there was an error opening or reading a file, 0 otherwise
# TYPE node_textfile_scrape_error gauge
node_textfile_scrape_error 0
```
