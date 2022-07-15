---
title: PromQL Parser
---

This tool enables you to validate the format of PromQL queries. It also produces
a prettified query.

## Usage

To use this tool, simply paste the content of your query in the
following text area.
Then, click on the "Parse" button.

You can click on the following button to load a PromQL query:

{{< unsafe >}}
<button onClick="loadexample();" id="exampleButton" disabled>Load sample</button>
{{< /unsafe >}}

## Security and privacy

The input is parsed in your browser and is not sent to our servers. This tool is
based on the official
[prometheus](https://pkg.go.dev/github.com/prometheus/prometheus@main/promql/parser#Prettify) library and is
cross compiled to [WASM](https://webassembly.org/), so that it runs natively in
your browser.

## PromQL input

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
WebAssembly.instantiateStreaming(fetch("/promqlparser.wasm"),
        go.importObject).then((result) => {
           go.run(result.instance);
});

</script>
<textarea style="display:block; width: 100%; height: 20em; margin-bottom: 2em;" name="promqlInput" id="promqlInput"></textarea>
<button onClick="parsepromql(promqlInput.value);" id="runButton" disabled>Parse</button>
<div id="resultDiv"></div>
{{< /unsafe >}}

