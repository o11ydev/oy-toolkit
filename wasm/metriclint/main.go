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
	"strings"
	"syscall/js"

	"github.com/prometheus/client_golang/prometheus/testutil/promlint"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("metriclint", js.FuncOf(metriclint))
	jsDoc := js.Global().Get("document")
	btn := jsDoc.Call("getElementById", "runButton")
	btn.Set("disabled", false)

	<-c
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
	  <strong>Error parsing metrics.</strong><br>
	  %s
	  </blockquote>
	`, err.Error()))
		return nil
	}
	if len(problems) == 0 {
		res.Set("innerHTML", `
	<blockquote class="gdoc-hint tip">
	  <strong>Success.</strong><br>
	  </blockquote>
	`)
		return nil
	}

	var pbs string

	for _, p := range problems {
		pbs += fmt.Sprintf("<li><code>%s</code>: %s</li>", p.Metric, p.Text)
	}

	res.Set("innerHTML", fmt.Sprintf(`
	<blockquote class="gdoc-hint warning">
	  <strong>Warnings found.</strong><br>
	  <ul>%s</ul>
	  </blockquote>
	`, pbs))

	return nil
}
