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

package client

import (
	"fmt"
	"os"

	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/common/config"
	"github.com/alecthomas/kingpin/v2"
	"gopkg.in/yaml.v2"
)

type Config struct {
	url  string
	file string
}

func InitCliFlags() *Config {
	var c Config
	kingpin.Flag("prometheus.url", "URL of the Prometheus server.").Default("http://127.0.0.1:9090").StringVar(&c.url)
	kingpin.Flag("client.config", "Path to a HTTP client configuration.").StringVar(&c.file)
	return &c
}

func NewClient(c *Config) (api.Client, error) {
	cfg := config.DefaultHTTPClientConfig
	if c.file != "" {
		dat, err := os.ReadFile(c.file)
		if err != nil {
			return nil, fmt.Errorf("error reading client config %s: %w", c.file, err)
		}
		err = yaml.UnmarshalStrict(dat, &cfg)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling client config %s: %w", c.file, err)
		}
	}
	rt, err := config.NewRoundTripperFromConfig(cfg, "oy-toolkit")
	if err != nil {
		return nil, fmt.Errorf("error creating roundtripper: %w", err)
	}
	return api.NewClient(api.Config{
		Address:      c.url,
		RoundTripper: rt,
	})
}
