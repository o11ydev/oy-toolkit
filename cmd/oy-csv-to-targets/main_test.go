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
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

var mainPath = os.Args[0]

func TestMain(m *testing.M) {
	for i, arg := range os.Args {
		if arg == "-test.main" {
			os.Args = append(os.Args[:i], os.Args[i+1:]...)
			main()
			return
		}
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestCSVToTargets(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	run := exec.Command(mainPath, "-test.main", "--input.file=testdata/targets.csv", "--output.file=testdata/targets.json")

	// Log stderr in case of failure.
	stderr, err := run.StderrPipe()
	require.NoError(t, err)
	go func() {
		slurp, _ := ioutil.ReadAll(stderr)
		t.Log(string(slurp))
	}()

	err = run.Run()
	require.NoError(t, err)

	compareJSON(t, "testdata/expected_targets.json", "testdata/targets.json")
}

func compareJSON(t *testing.T, file1, file2 string) {
	data1, err := ioutil.ReadFile(file1)
	require.NoError(t, err)

	data2, err := ioutil.ReadFile(file2)
	require.NoError(t, err)

	var obj1 interface{}
	err = json.Unmarshal(data1, &obj1)
	require.NoError(t, err)

	var obj2 interface{}
	err = json.Unmarshal(data2, &obj2)
	require.NoError(t, err)

	require.Equal(t, obj1, obj2)
}
