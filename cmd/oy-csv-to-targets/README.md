# oy-csv-to-targets

*oy-csv-to-targets* takes a list of Prometheus targets as a CSV file as input
and produces a JSON file.

The resulting JSON file can then be used by the `file_sd` discovery mechanism to
dynamically discover and scrape the Prometheus targets. This allows the targets
to be easily managed and updated without having to manually update the
Prometheus configuration file, if your team is not fluent with JSON.

If no output file is specified, the output goes to stdout.

## Example usage

With the following CSV as an input:

```
,datacenter,availability_zone
prometheus1:9090,dc1,az1
prometheus2:9090,dc1,az2
prometheus3:9090,dc2,az1
```

The following command:

```
./oy-csv-to-targets --input.file targets.csv --output.file targets.json
```

Would produce the following `targets.json` file:

```json
[
    {
        "labels": {
            "availability_zone": "az1",
            "datacenter": "dc1"
        },
        "targets": [
            "prometheus1:9090"
        ]
    },
    {
        "labels": {
            "availability_zone": "az2",
            "datacenter": "dc1"
        },
        "targets": [
            "prometheus2:9090"
        ]
    },
    {
        "labels": {
            "availability_zone": "az1",
            "datacenter": "dc2"
        },
        "targets": [
            "prometheus3:9090"
        ]
    }
]
```

You can then configure your Prometheus as follows

```
scrape_configs:
- job_name: prometheus
  file_sd_configs:
    - files: [targets.json]
```
