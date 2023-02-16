// Copyright The o11y toolkit Authors
// spdx-license-identifier: apache-2.0
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at:
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

//go:build js && wasm
// +build js,wasm

package main

import (
	"fmt"
	"strings"
	"syscall/js"

	"github.com/prometheus/client_golang/prometheus/testutil/promlint"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("metricslint", js.FuncOf(metriclint))
	js.Global().Set("loadexample", js.FuncOf(loadexample))
	jsDoc := js.Global().Get("document")
	jsDoc.Call("getElementById", "runButton").Set("disabled", false)
	jsDoc.Call("getElementById", "exampleButton").Set("disabled", false)
	jsDoc.Call("getElementById", "loadingWarning").Get("style").Set("display", "none")

	<-c
}

func loadexample(this js.Value, args []js.Value) interface{} {
	jsDoc := js.Global().Get("document")
	res := jsDoc.Call("getElementById", "metricInput")
	res.Set("value", `# HELP process_cpu_seconds_total Total user and system CPU time spent in seconds.
# TYPE process_cpu_seconds_total counter
process_cpu_seconds_total 203674.05
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds 1024
# HELP process_open_fds Number of open file descriptors.
# TYPE process_open_fds gauge
process_open_fds 10
`)
	return nil
}

func metriclint(this js.Value, args []js.Value) interface{} {
	jsDoc := js.Global().Get("document")
	res := jsDoc.Call("getElementById", "resultDiv")

	metrics := strings.NewReader(args[0].String() + "\n")
	l := promlint.New(metrics)
	problems, err := l.Lint()
	if err != nil {
		res.Set("innerHTML", fmt.Sprintf(`
	<blockquote class="gdoc-hint danger">
	<div class="gdoc-hint__title flex align-center">
		<svg class="gdoc-icon gdoc_dangerous"><use xlink:href="#gdoc_dangerous"></use></svg>
		<span>Parsing error</span>
	</div>
	<div class="gdoc-hint__text">
	  %s
	</div>
	</blockquote>
	`, err.Error()))
		return nil
	}

	if strings.TrimSpace(args[0].String()) == "" {
		res.Set("innerHTML", `
	<blockquote class="gdoc-hint important">
	<div class="gdoc-hint__title flex align-center">
		<svg class="gdoc-icon gdoc_error_outline"><use xlink:href="#gdoc_error_outline"></use></svg>
		<span>No input</span>
	</div>
	<div class="gdoc-hint__text">
	  The input provided is empty. Please paste metrics into the text area.
	</div>
	</blockquote>
	`)
		return nil
	}

	if len(problems) == 0 {
		res.Set("innerHTML", `
	<blockquote class="gdoc-hint tip">
	<div class="gdoc-hint__title flex align-center">
		<svg class="gdoc-icon gdoc_check"><use xlink:href="#gdoc_check"></use></svg>
		<span>Success</span>
	</div>
	<div class="gdoc-hint__text">
	  Input has been parsed successfully.
	</div>
	</blockquote>
	`)
		return nil
	}

	var pbs string

	for _, p := range problems {
		pbs += fmt.Sprintf("<li><code>%s</code>: %s</li>", p.Metric, p.Text)
	}

	res.Set("innerHTML", fmt.Sprintf(`
	<blockquote class="gdoc-hint important">
	<div class="gdoc-hint__title flex align-center">
		<svg class="gdoc-icon gdoc_fire"><use xlink:href="#gdoc_fire"></use></svg>
		<span>Issues found</span>
	</div>
	<div class="gdoc-hint__text">
	  The input can be parsed but there are linting issues:
	  <ul>%s</ul>
	</div>
	</blockquote>
	`, pbs))

	return nil
}
