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

package cmd

import (
	"os"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/common/promlog/flag"
	"github.com/prometheus/common/version"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func InitCmd(name string) log.Logger {
	promlogConfig := &promlog.Config{}
	flag.AddFlags(kingpin.CommandLine, promlogConfig)
	kingpin.Version(version.Print(name))
	kingpin.CommandLine.UsageWriter(os.Stdout)
	kingpin.HelpFlag.Short('h')
	kingpin.Parse()
	return newLogger(promlogConfig)
}

func newLogger(config *promlog.Config) log.Logger {
	var l log.Logger
	if config.Format != nil && config.Format.String() == "json" {
		l = log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	} else {
		l = log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))
	}

	if config.Level != nil {
		var lvl level.Option
		switch config.Level.String() {
		case "debug":
			lvl = level.AllowDebug()
		case "info":
			lvl = level.AllowInfo()
		case "warn":
			lvl = level.AllowWarn()
		case "error":
			lvl = level.AllowError()
		default:
			lvl = level.AllowDebug()
		}
		l = level.NewFilter(l, lvl)
	}
	return l
}
