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
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"

	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"

	"github.com/o11ydev/oy-toolkit/util/period"
)

// newGroupCollector returns a new Collector exposing metrics from rules
// evaluated at periodic intervals.
func newGroupsCollector(logger log.Logger, client api.Client, groups []Group) (prometheus.Collector, error) {
	c := &groupsCollector{
		groups: groups,
		logger: logger,
		client: client,
		cache:  NewCache(24 * time.Hour),
	}
	return c, nil
}

type groupsCollector struct {
	groups []Group
	logger log.Logger
	client api.Client
	cache  *Cache
}

// Collect implements the Collector interface.
func (c *groupsCollector) Collect(ch chan<- prometheus.Metric) {
	success := func(val float64) {
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc("periodic_rule_queries_success", "", nil, nil),
			prometheus.GaugeValue,
			val,
		)
	}

	metrics := []prometheus.Metric{}
	now := time.Now().Local()
	for _, g := range c.groups {
		for _, p := range period.CalculatePeriods(now, g.TimePeriod, time.Duration(g.Lookback), g.IncludeIncompleteRanges) {
			for _, r := range g.Rules {
				m, err := QueryRangeToMetric(c.logger, c.cache, p, r.Expr, c.client, r.Record, r.Labels)
				if err != nil {
					success(0)
					return
				}
				metrics = append(metrics, m...)
			}
		}
	}
	for _, metric := range metrics {
		ch <- metric
	}
	success(1)
}

func (c *groupsCollector) Describe(ch chan<- *prometheus.Desc) {
}

func QueryRangeToMetric(logger log.Logger, c *Cache, p period.Period, query string, client api.Client, metricName string, additionalLabels map[string]string) ([]prometheus.Metric, error) {
	queryAPI := v1.NewAPI(client)

	pq := replaceForPeriod(p, query)

	var (
		result model.Value
		err    error
		ok     bool
	)

	var found bool
	key := fmt.Sprintf("%d/%s", p.End.Unix(), pq)
	if p.Complete {
		if result, ok = c.Get(key); ok {
			found = true
		}
	}
	if !found {
		if pq != query {
			level.Debug(logger).Log("msg", "replaced query", "period_start", p.Start, "period_end", p.End, "period_complete", p.Complete, "query", query, "new_query", pq)
		}
		result, _, err = queryAPI.Query(context.Background(), pq, p.End)
		if err != nil {
			level.Error(logger).Log("err", err)
			return nil, err
		}
		if p.Complete {
			c.Set(key, result)
		}
	}

	var metrics []prometheus.Metric
	metricLabels := make(prometheus.Labels, len(additionalLabels))
	for k, v := range additionalLabels {
		metricLabels[k] = replaceForPeriod(p, v)
	}

	for _, vec := range result.(model.Vector) {
		labels := make(prometheus.Labels, len(vec.Metric)+len(metricLabels))
		for k, v := range vec.Metric {
			if k == model.MetricNameLabel {
				continue
			}
			labels[string(k)] = string(v)
		}
		for k, v := range metricLabels {
			if k == model.MetricNameLabel {
				continue
			}
			labels[k] = v
		}

		sample := prometheus.MustNewConstMetric(
			prometheus.NewDesc(metricName, "", nil, labels),
			prometheus.UntypedValue,
			float64(vec.Value),
		)

		metrics = append(metrics, sample)
	}

	return metrics, nil
}

func replaceForPeriod(p period.Period, v string) string {
	d := p.End.Sub(p.Start)
	replacements := map[string]string{
		"RANGE": fmt.Sprintf("%.0fs", d.Seconds()),
	}
	for i := 0; i <= 100; i++ {
		replacements[strconv.Itoa(i)+"_PC_TIMESTAMP"] = fmt.Sprintf("%d", p.Start.Add(time.Duration(i/100)*d).Unix())
		replacements[strconv.Itoa(i)+"_PC_RANGE"] = fmt.Sprintf("%.0fs", float64(i)*d.Seconds()/100)
	}

	return os.Expand(v, func(s string) string {
		if s == "$" {
			return "$"
		}
		if v, ok := replacements[s]; ok {
			return v
		}
		if strings.HasPrefix(s, "_") {
			return p.Start.Format(strings.TrimPrefix(s, "-"))
		}
		return p.End.Format(s)
	})
}
