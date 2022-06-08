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
	"errors"
	"fmt"
	"strconv"
	"syscall/js"

	"golang.org/x/crypto/bcrypt"
	"gopkg.in/yaml.v2"
)

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("generateUsers", js.FuncOf(pwgen))
	jsDoc := js.Global().Get("document")
	jsDoc.Call("getElementById", "runButton").Set("disabled", false)
	jsDoc.Call("getElementById", "loadingWarning").Get("style").Set("display", "none")
	js.Global().Call("addUser")

	<-c
}

func pwgen(this js.Value, args []js.Value) interface{} {
	jsDoc := js.Global().Get("document")
	res := jsDoc.Call("getElementById", "resultDiv")
	users := jsDoc.Call("querySelectorAll", `[name="username"]`)
	passwords := jsDoc.Call("querySelectorAll", `[name="password"]`)
	cost := jsDoc.Call("querySelectorAll", `[name="cost"]`)
	c, err := strconv.ParseInt(cost.Index(0).Get("value").String(), 10, 32)
	if err != nil {
		mkErr(err)
		return nil
	}
	userspw := make(map[string]string, users.Length())
	for i := 0; i < users.Length(); i++ {
		user := users.Index(i).Get("value").String()
		pw := passwords.Index(i).Get("value").String()
		if user == "" {
			mkErr(errors.New("username can't be empty"))
			return nil
		}
		if pw == "" {
			mkErr(fmt.Errorf("password for %q can't be empty", user))
			return nil
		}
		if _, ok := userspw[user]; ok {
			mkErr(fmt.Errorf("duplicate user %q", user))
			return nil
		}
		gpw, err := bcrypt.GenerateFromPassword([]byte(pw), int(c))
		if err != nil {
			mkErr(err)
			return nil
		}
		userspw[user] = string(gpw)
	}

	out, err := yaml.Marshal(struct {
		Users map[string]string `yaml:"basic_auth_users"`
	}{Users: userspw})
	if err != nil {
		mkErr(err)
		return nil
	}

	res.Set("innerHTML", fmt.Sprintf(`
	<h2>Web configuration file</h2>
	<p>Write the following content as <code>web.yml</code> file and start Prometheus with <code>--web.config.file=web.yml</code></p>
	<pre class="chroma"><code class="language-yaml" data-lang="yaml">%s</code></pre>
	`, out))
	return nil
}

func mkErr(err error) {
	jsDoc := js.Global().Get("document")
	res := jsDoc.Call("getElementById", "resultDiv")
	res.Set("innerHTML", fmt.Sprintf(`
	<blockquote class="gdoc-hint danger">
	<div class="gdoc-hint__title flex align-center">
		<svg class="gdoc-icon gdoc_dangerous"><use xlink:href="#gdoc_dangerous"></use></svg>
		<span>Error</span>
	</div>
	<div class="gdoc-hint__text">
	  %s
	</div>
	</blockquote>
	`, err.Error()))
}
