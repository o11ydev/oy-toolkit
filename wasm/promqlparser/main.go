// Copyright 2022 The o11y toolkit Authors
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
	"syscall/js"

	"github.com/prometheus/prometheus/promql/parser"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("parsepromql", js.FuncOf(metriclint))
	js.Global().Set("loadexample", js.FuncOf(loadexample))
	jsDoc := js.Global().Get("document")
	jsDoc.Call("getElementById", "runButton").Set("disabled", false)
	jsDoc.Call("getElementById", "exampleButton").Set("disabled", false)
	jsDoc.Call("getElementById", "loadingWarning").Get("style").Set("display", "none")

	<-c
}

func loadexample(this js.Value, args []js.Value) interface{} {
	jsDoc := js.Global().Get("document")
	res := jsDoc.Call("getElementById", "promqlInput")
	res.Set("value", `rate(foo[5s])`)
	return nil
}

func metriclint(this js.Value, args []js.Value) interface{} {
	jsDoc := js.Global().Get("document")
	res := jsDoc.Call("getElementById", "resultDiv")

	promql := args[0].String()
	_, err := parser.ParseExpr(promql)
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
