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
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/alecthomas/kingpin.v2"

	promlabels "github.com/prometheus/prometheus/model/labels"

	"github.com/o11ydev/oy-toolkit/util/cmd"
)

var (
	input  = kingpin.Flag("input.file", "Path to a CSV file to use as an input.").PlaceHolder("input.csv").Required().String()
	output = kingpin.Flag("output.file", "Path to a json file to use as an output.").PlaceHolder("targets.json").String()
)

type target struct {
	Labels  promlabels.Labels `json:"labels"`
	Targets []string          `json:"targets"`
}

func main() {
	logger := cmd.InitCmd("oy-csv-to-targets")

	err := csvToJSON(*input, *output)
	if err != nil {
		logger.Log("error", err.Error())
	}
}

// Convert a CSV file to a JSON with Prometheus targets.
func csvToJSON(csvFile, jsonFile string) error {
	f, err := os.Open(csvFile)
	if err != nil {
		return fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	// Read the first line of the CSV to get the headers (i.e. label names).
	headers, err := r.Read()
	if err != nil {
		return fmt.Errorf("failed to read headers from CSV: %w", err)
	}

	targets := []target{}

	// Read the rest of the CSV to get the target values
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read record from CSV: %w", err)
		}

		// Create a map to hold the label values
		labels := map[string]string{}

		// Loop through the headers and record values to create the label
		// values.
		// Ignore the first one.
		for i, header := range headers {
			if i > 0 {
				labels[header] = record[i]
			}
		}

		// Create a new Prometheus target using the label values and the first
		// record value as the target address.
		t := target{
			Labels:  promlabels.FromMap(labels),
			Targets: []string{record[0]},
		}

		targets = append(targets, t)
	}

	b, err := json.Marshal(targets)
	if err != nil {
		return fmt.Errorf("failed to marshal targets into JSON: %w", err)
	}

	// If there is no file, pretty print to stdout.
	if jsonFile == "" {
		var out bytes.Buffer
		err = json.Indent(&out, b, "", "    ")
		if err != nil {
			return fmt.Errorf("failed to pretty print JSON: %w", err)
		}

		fmt.Println(out.String())
		return nil
	}

	// Otherwise, output to the file using a temp file.
	// Create a temp file in the same directory as the destination file.
	tmpFile, err := ioutil.TempFile(filepath.Dir(jsonFile), "tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	_, err = tmpFile.Write(b)
	if err != nil {
		tmpFile.Close()
		os.Remove(tmpFile.Name())
		return fmt.Errorf("failed to write to temp file: %w", err)
	}

	err = tmpFile.Close()
	if err != nil {
		os.Remove(tmpFile.Name())
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Rename the temp file to the destination file.
	err = os.Rename(tmpFile.Name(), jsonFile)
	if err != nil {
		os.Remove(tmpFile.Name())
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	return nil
}
