# oy-periodic-queries

*oy-periodic-queries* is a tool that allows you to evaluate Prometheus
recording rules and export the results as metrics, with defined boundaries such
as monthly or weekly (*Only monthly is implemented at the moment*).

It uses calendar months.

With this exporter, you can easily calculate monthly resource usage, or any
other metric that requires periodic evaluation. Simply specify the recording
rule and time boundaries, and the exporter will run the rule and provide you
with accurate, timely metrics.

In the recording rule and labels, you have access to the following variables, that will be
replaced:

| Variable | Meaning |
|----------|---------|
| $RANGE | The range of the query, in duration |
| $0_PC_TIMESTAMP | The start time of the query, in unix timestamp |
| $x_PC_TIMESTAMP | The time located a x % of the range, in unix timestamp, x
between 0 and 100 |
| $100_PC_TIMESTAMP | The end time of the query, in unix timestamp |
| $x_PC_RANGE | A duration that represents x % of the current range, x between 0
and 100 |

Other values in `${}` are formatted with [Go's
Time.Format](https://pkg.go.dev/time#Time.Format), e.g. `"${2006-01}"` is turned
into the month and day of the end of the range. Prepend `_` to have it relative
to the start of the range: `"${_2006-01}"`.

## Configuration

The syntax of a rule file is:

```
groups:
  [ - <rule_group> ]
```

### `<rule_group>`

```
# The name of the group. Must be unique within a file.
name: <string>

# The time period covered by the metrics.
time_period: [ <string> | default = "monthly" ]

# The maximum time to look back.
lookback: [ <duration> | default = "12w" ]

rules:
  [ - <rule> ... ]
```

### `<rule>`

The syntax for recording rules is:

```
# The name of the time series to output to. Must be a valid metric name.
record: <string>

# The PromQL expression to evaluate. Every evaluation cycle this is
# evaluated at the current time, and the result recorded as a new set of
# time series with the metric name as given by 'record'.
expr: <string>

# Labels to add or overwrite before storing the result.
labels:
  [ <labelname>: <labelvalue> ]
```

## Example rule file

A simple example rules file would be:

```
groups:
- name: prometheus network usage
  time_period: monthly
  lookback: 365d
  include_incomplete_ranges: true
  rules:
    - expr: |
        increase(node_network_receive_bytes_total{instance=~"prometheus.*"}[$RANGE])
        and last_over_time(node_network_receive_bytes_total{instance=~"prometheus.*"}[${1_PC_RANGE}] @ ${1_PC_TIMESTAMP})
        and last_over_time(node_network_receive_bytes_total{instance=~"prometheus.*"}[${1_PC_RANGE}] @ ${100_PC_TIMESTAMP})
      record: prometheus:node_network_receive_bytes_total:monthly
      labels:
        month: "${2006-01}"
    - expr: |
        increase(node_network_transmit_bytes_total{instance=~"prometheus.*"}[$RANGE])
        and last_over_time(node_network_transmit_bytes_total{instance=~"prometheus.*"}[${1_PC_RANGE}] @ ${1_PC_TIMESTAMP})
        and last_over_time(node_network_transmit_bytes_total{instance=~"prometheus.*"}[${1_PC_RANGE}] @ ${100_PC_TIMESTAMP})
      record: prometheus:node_network_transmit_bytes_total:monthly
      labels:
        month: "${2006-01}"
```


In the queries, the following part is used to only select the metrics that
existed in the beginning and the end of the different months.
```
        and last_over_time(node_network_transmit_bytes_total{instance=~"prometheus.*"}[${1_PC_RANGE}] @ ${1_PC_TIMESTAMP})
        and last_over_time(node_network_transmit_bytes_total{instance=~"prometheus.*"}[${1_PC_RANGE}] @ ${100_PC_TIMESTAMP})
```

Output of /metrics, when launched with `--web.disable-exporter-metrics`:


```
# HELP periodic_rule_queries_success 
# TYPE periodic_rule_queries_success gauge
periodic_rule_queries_success 1
# HELP prometheus:node_network_receive_bytes_total:monthly 
# TYPE prometheus:node_network_receive_bytes_total:monthly untyped
prometheus:node_network_receive_bytes_total:monthly{device="ens3",environment="o11ylab",instance="prometheus02.example.com",job="node",month="2023-01"} 2.2096172024710458e+11
prometheus:node_network_receive_bytes_total:monthly{device="ens3",environment="o11ylab",instance="prometheus11.example.com",job="node",month="2023-01"} 1.4455463702554254e+11
prometheus:node_network_receive_bytes_total:monthly{device="lo",environment="o11ylab",instance="prometheus02.example.com",job="node",month="2023-01"} 9.635369117702129e+08
prometheus:node_network_receive_bytes_total:monthly{device="lo",environment="o11ylab",instance="prometheus11.example.com",job="node",month="2023-01"} 2.4574080228794675e+09
# HELP prometheus:node_network_transmit_bytes_total:monthly 
# TYPE prometheus:node_network_transmit_bytes_total:monthly untyped
prometheus:node_network_transmit_bytes_total:monthly{device="ens3",environment="o11ylab",instance="prometheus02.example.com",job="node",month="2023-01"} 1.8763866363794147e+11
prometheus:node_network_transmit_bytes_total:monthly{device="ens3",environment="o11ylab",instance="prometheus11.example.com",job="node",month="2023-01"} 3.813239344424051e+11
prometheus:node_network_transmit_bytes_total:monthly{device="lo",environment="o11ylab",instance="prometheus02.example.com",job="node",month="2023-01"} 9.635369117702129e+08
prometheus:node_network_transmit_bytes_total:monthly{device="lo",environment="o11ylab",instance="prometheus11.example.com",job="node",month="2023-01"} 2.4574080228794675e+09
```
