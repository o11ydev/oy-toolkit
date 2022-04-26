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

package client

import (
	"github.com/prometheus/client_golang/api"
	"gopkg.in/alecthomas/kingpin.v2"
)

type Config struct {
	url string
}

func InitCliFlags() *Config {
	var c Config
	kingpin.Flag("prometheus.url", "URL of the Prometheus server.").Default("http://127.0.0.1:9090").StringVar(&c.url)
	return &c
}

func NewClient(c *Config) (api.Client, error) {
	return api.NewClient(api.Config{
		Address: c.url,
	})
}
