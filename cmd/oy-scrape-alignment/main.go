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

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/prometheus/client_golang/api"
	apiv1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
	"gopkg.in/alecthomas/kingpin.v2"

	"github.com/o11ydev/oy-toolkit/util/client"
	"github.com/o11ydev/oy-toolkit/util/cmd"
)

var (
	metric        = kingpin.Flag("metric", "Metric to use to determine alignment.").Default("up").String()
	lookback      = kingpin.Flag("lookback", "How much time to look in the past for scrapes.").Default("1h").Duration()
	divisor       = kingpin.Flag("divisor", "Divisor to use to determine if a scrape is aligned.").Default("1s").Duration()
	unalignedOnly = kingpin.Flag("log.unaligned-only", "Only log unaligned targets.").Bool()
	quiet         = kingpin.Flag("log.results-only", "Only log final result.").Bool()
)

func main() {
	c := client.InitCliFlags()
	logger := cmd.InitCmd("oy-scrape-alignment")

	promClient, err := client.NewClient(c)
	if err != nil {
		level.Error(logger).Log("msg", "Can't create Prometheus client", "err", err)
		os.Exit(1)
	}

	analyzeScrapeAlignment(logger, promClient)
}

func analyzeScrapeAlignment(logger log.Logger, promClient api.Client) {
	api := apiv1.NewAPI(promClient)
	v, warnings, err := api.Query(context.Background(), fmt.Sprintf("%s[%dms]", *metric, lookback.Milliseconds()), time.Now())
	if err != nil {
		level.Error(logger).Log("msg", "Can't query up metrics", "err", err)
		os.Exit(1)
	}
	for w := range warnings {
		if err != nil {
			level.Warn(logger).Log("msg", w)
		}
	}

	if v.Type() != model.ValMatrix {
		if err != nil {
			level.Error(logger).Log("msg", "Wrong return type", "expected", model.ValMatrix, "got", v.Type())
			os.Exit(1)
		}
	}

	result, _ := v.(model.Matrix)
	var goodTargets, badTargets int
	var maxTarget int64
	for _, r := range result {
		var good, bad float64
		var max int64
		var lastTs time.Time
		for _, s := range r.Values {
			if lastTs.IsZero() {
				lastTs = s.Timestamp.Time()
				continue
			}
			diff := s.Timestamp.Time().Sub(lastTs).Milliseconds()
			ok := diff%divisor.Milliseconds() == 0
			level.Debug(logger).Log("metric", r.Metric.String(), "prev", lastTs, "current", s.Timestamp.Time(), "difference", diff, "aligned", ok)
			if ok {
				good++
			} else {
				bad++
				if diff%divisor.Milliseconds() > divisor.Milliseconds()/2 {
					diff = divisor.Milliseconds() - diff%divisor.Milliseconds()
				}
				if diff > max {
					max = diff % divisor.Milliseconds()
					if max > maxTarget {
						maxTarget = max
					}
				}
			}
			lastTs = s.Timestamp.Time()
		}

		if (bad != 0 || !*unalignedOnly) && !*quiet {
			level.Info(logger).Log("metric", r.Metric.String(), "aligned", good, "unaligned", bad, "max_ms", max, "pc", fmt.Sprintf("%.2f%%", 100*good/(bad+good)))
		}

		if bad == 0 {
			goodTargets++
		} else {
			badTargets++
		}
	}
	level.Info(logger).Log("msg", "overall results", "aligned_targets", goodTargets, "unaligned_targets", badTargets, "max_ms", maxTarget)
}
