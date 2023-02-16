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
	"fmt"
	"os"

	"github.com/prometheus/common/model"
	"gopkg.in/yaml.v2"
)

type groupfile struct {
	Groups []Group `yaml:"groups"`
}

type Group struct {
	Name                    string         `yaml:"name"`
	TimePeriod              string         `yaml:"time_period"`
	Lookback                model.Duration `yaml:"lookback"`
	Rules                   []Rule         `yaml:"rules"`
	IncludeIncompleteRanges bool           `yaml:"include_incomplete_ranges"`
}

type Rule struct {
	Record string            `yaml:"record"`
	Expr   string            `yaml:"expr"`
	Labels map[string]string `yaml:"labels"`
}

func (g *Group) Validate() error {
	if g.TimePeriod != "monthly" {
		return fmt.Errorf("Only monthly period is supported, got %q", g.TimePeriod)
	}
	return nil
}

func (g *Group) UnmarshalYAML(unmarshal func(interface{}) error) error {
	*g = Group{
		IncludeIncompleteRanges: true,
	}
	type plain Group
	if err := unmarshal((*plain)(g)); err != nil {
		return err
	}

	return g.Validate()
}

func loadFiles(files []string) ([]Group, error) {
	g := make([]Group, 0)
	for _, file := range files {
		data, err := os.ReadFile(file)
		if err != nil {
			return nil, fmt.Errorf("error reading rule file %s: %w", file, err)
		}
		gf := groupfile{}
		err = yaml.UnmarshalStrict(data, &gf)
		if err != nil {
			return nil, fmt.Errorf("error unmarshaling rule file %s: %w", file, err)
		}
		for _, group := range gf.Groups {
			if err := group.Validate(); err != nil {
				return nil, err
			}
			g = append(g, group)
		}
	}
	return g, nil
}
