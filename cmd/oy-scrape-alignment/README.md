# oy-scrape-alignment

*oy-scrape-alignment* queries a Prometheus server to see how aligned scrapes are.
Perfect scrape alignment happens when the distance between all the scrapes are
exactly the same. It enables Prometheus [delta-of-delta](https://github.com/prometheus/prometheus/blob/main/tsdb/docs/bstream.md)
encoding to reduce significantly the size of the blocks.

## Example usage

```shell
$ ./oy-scrape-alignment --prometheus.url=https://prometheus.demo.do.prometheus.io/ --log.results-only
ts=2022-04-26T12:11:34.659Z caller=main.go:117 level=info msg="overall results" aligned_targets=0 unaligned_targets=10 max_ms=25
```

This means that the maximum deviation seen in your scrape jobs is 25ms. You
could set `--scrape.timestamp-tolerance=25ms` to reduce your disk usage over
time, by enabling Prometheus to correct timestamps up to 25ms.


## Prometheus limits

Prometheus will only apply timestamp tolerance up to 1%. If your scrape interval
is 30s, you can only adjust timestamps up to 300ms. Setting a 500ms tolerance
will have no effects on jobs with a scrape interval lower than 500s, even if the
deviation is tiny.

## Plotting the output

By using `--plot.file=scrape.png`, you can generate a PNG file which shows the
scrape (mis-)alignment with an histogram.
