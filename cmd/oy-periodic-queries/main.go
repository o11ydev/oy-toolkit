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

	"github.com/o11ydev/oy-toolkit/util/client"
	"github.com/o11ydev/oy-toolkit/util/cmd"
	"github.com/o11ydev/oy-toolkit/util/http"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	ruleFiles := kingpin.Arg("rule-file", "File to read periodic recording rules from.").Required().Strings()
	c := client.InitCliFlags()
	logger := cmd.InitCmd("oy-periodic-files")

	promClient, err := client.NewClient(c)
	if err != nil {
		level.Error(logger).Log("msg", "Can't create Prometheus client", "err", err)
		os.Exit(1)
	}

	groups, err := loadFiles(*ruleFiles)
	if err != nil {
		level.Error(logger).Log("err", err)
		os.Exit(1)
	}
	collector, err := newGroupsCollector(logger, promClient, groups)
	if err != nil {
		level.Error(logger).Log("err", err)
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
