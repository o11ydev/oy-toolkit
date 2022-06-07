---
title: Prometheus Password Generator
---

This form enables you to generate a [Prometheus web.yml
file](https://prometheus.io/docs/prometheus/latest/configuration/https/) to
secure your Prometheus endpoints with basic authentication.

To do this, you need to hash passwords with [bcrypt](https://en.wikipedia.org/wiki/Bcrypt).
This tool hashes the passwords directly in your browser, in such a way that we
do not receive the passwords you are generating.

Once the file is generated, you can optionally append your TLS server
configuration to the file, then start Prometheus with `--web.config.file`
pointing to your newly created file.

This file is also compatible with Alertmanager, Pushgateway, Node Exporter and
other official exporters.

## How to

Enter the usernames and the password, then press the generate button to generate
the file.

You can add and remove users with the `Remove` and `Add user` buttons.


## Security and privacy

The input is parsed in your browser and is not sent to our servers. This tool is
cross compiled to [WASM](https://webassembly.org/), so that it runs natively in
your browser.

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
WebAssembly.instantiateStreaming(fetch("/pwgen.wasm"),
        go.importObject).then((result) => {
           go.run(result.instance);
});

addUser = function(){
    tb = document.getElementById('usersTable');
    newRow = tb.insertRow();
     newRow.insertCell().innerHTML='<input type="text" name="username" placeholder="username">';
     newRow.insertCell().innerHTML='<input type="password" name="password" placeholder="password">';
     newRow.insertCell().innerHTML='<input type="button" value="Remove" onclick="removeUser(this)">';
};

switchViz = function(t){
    pw = document.querySelectorAll('[name="password"]');
    for (i = 0; i < pw.length; ++i) {
        if (pw[i].type === "password") {
            t.innerHTML="Hide passwords";
            pw[i].type = "text";
        } else {
            t.innerHTML="Show passwords";
            pw[i].type = "password";
        }
    }
};

removeUser = function(t) {
    var p = t.parentNode.parentNode;
    p.parentNode.removeChild(p);
};

</script>

<table id="usersTable">
<tr>
<th>Username</th>
<th>Password</th>
<th></th>
</tr>
</table>

<button onClick="addUser();" id="addUserButton">New user</button>
<button onClick="switchViz(this);" id="switchViz">Show passwords</button>
<button onClick="generateUsers();" id="runButton" disabled>Generate</button>
<div id="resultDiv"></div>
{{< /unsafe >}}

