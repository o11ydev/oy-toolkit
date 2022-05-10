---
title: /metrics lint
---

This tool enables you to validate the format of Prometheus metrics, and make
sure they can be scraped by a Prometheus server.

Prometheus supports two exposition formats: the Prometheus text-based exposition
format and [OpenMetrics](https://openmetrics.io). The text-based exposition
format is widespread, and many applications and client libraries supports it.
Additionally, it can be used from scripts, to push metrics to the
[Pushgateway](https://github.com/prometheus/pushgateway) or written to `*.prom` files, to
be collected by the [textfile collectors](https://github.com/prometheus/node_exporter#textfile-collector)
(available in both the [Node Exporter](https://github.com/prometheus/node_exporter) and the
[Windows Exporter](https://github.com/prometheus-community/windows_exporter)).

Our toolkit also provides [oy-expose](/oy-expose), a standalone tool that can
expose the metrics of a file to be consumed by Prometheus.

## Usage

To use this tool, simply paste the content of your `*.prom` file, the body of
your Pushgateway request, or the output of a `/metrics` HTTP endpoint in the
following text area.
Then, click on the "Lint" button.

You can click on the following button to load a few metrics:

{{< unsafe >}}
<button onClick="loadexample();" id="exampleButton" disabled>Load sample</button>
{{< /unsafe >}}

## Security and privacy

The input is parsed in your browser and is not sent to our servers. This tool is
based on the official
[client_golang](https://github.com/prometheus/client_golang) library and is
cross compiled to [WASM](https://webassembly.org/), so that it runs natively in
your browser.

## Format specification

This tool uses the [Prometheus text-based exposition
format](https://prometheus.io/docs/instrumenting/exposition_formats/#exposition-formats).
OpenMetrics is not supported yet.

Everything is run locally from your browser, we do not receive or collect your
metrics.

## Command line tool

This utilise behaves like the `promtool check metrics` command, which is
downloadable with [Prometheus](https://prometheus.io/download).

## Metrics validation

{{< unsafe >}}
<div id="loadingWarning">
{{< /unsafe >}}

{{< hint type=caution title=Loading icon=gdoc_timer >}}
The application is loading. If this warning does not disappear, please make sure
that [your browser supports WASM](https://caniuse.com/wasm) and that javascript
is enabled.
{{< /hint >}}

{{< unsafe >}}
</div>
{{< /unsafe >}}

{{< unsafe >}}
<script src="/wasm_exec.js"></script>

<script>
if (!WebAssembly.instantiateStreaming) {
    // polyfill
    WebAssembly.instantiateStreaming = async (resp, importObject) => {
        const source = await (await resp).arrayBuffer();
        return await WebAssembly.instantiate(source, importObject);
    };
}

const go = new Go();
WebAssembly.instantiateStreaming(fetch("/metricslint.wasm"),
        go.importObject).then((result) => {
           go.run(result.instance);
});

</script>
<textarea style="display:block; width: 100%; height: 20em; margin-bottom: 2em;" name="metricInput" id="metricInput"></textarea>
<button onClick="metricslint(metricInput.value);" id="runButton" disabled>Lint</button>
<div id="resultDiv"></div>
{{< /unsafe >}}

