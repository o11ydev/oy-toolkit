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
	_ "embed"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

//go:embed query_result.json
var queryResult string

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

func TestExpose(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	testServer := httptest.NewServer(
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			require.Equal(t, "/api/v1/query", r.RequestURI)
			_, err := io.WriteString(w, queryResult)
			require.NoError(t, err)
		}),
	)
	t.Cleanup(testServer.Close)

	run := exec.Command(mainPath, "-test.main", "--prometheus.url="+testServer.URL)
	out, err := run.Output()
	t.Log(string(out))
	require.NoError(t, err)

	require.Contains(t, string(out), "level=info aligned_targets=4 unaligned_targets=6 max_ms=27")
}
