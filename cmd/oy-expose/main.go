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

package main

import (
	"context"
	"os"

	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/prometheus"

	"github.com/o11ydev/oy-toolkit/util/cmd"
	"github.com/o11ydev/oy-toolkit/util/collectors"
	"github.com/o11ydev/oy-toolkit/util/http"

	kingpin "github.com/alecthomas/kingpin/v2"
)

func main() {
	textFile := kingpin.Arg("metrics-file", "File to read metrics from.").Default("metrics").String()
	logger := cmd.InitCmd("oy-expose")

	collector, err := collectors.NewTextFileCollector(logger, *textFile)
	if err != nil {
		level.Error(logger).Log("msg", "Error creating file collector", "err", err)
		os.Exit(1)
	}

	r := prometheus.NewRegistry()
	r.MustRegister(collector)
	err = http.Serve(context.Background(), logger, r)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
}
