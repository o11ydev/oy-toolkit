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
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"os/exec"
	"testing"
	"time"

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

func TestExpose(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	testServer := httptest.NewServer(nil)
	testURL := testServer.URL
	testServer.Close()
	u, err := url.Parse(testURL)
	require.NoError(t, err)

	tmpfile, err := os.CreateTemp("", "metrics")
	require.NoError(t, err)
	defer os.Remove(tmpfile.Name())

	run := exec.Command(mainPath, "-test.main", "--web.listen-address="+u.Host, tmpfile.Name())

	// Log stderr in case of failure.
	stderr, err := run.StderrPipe()
	require.NoError(t, err)
	go func() {
		slurp, _ := ioutil.ReadAll(stderr)
		t.Log(string(slurp))
	}()

	err = run.Start()
	require.NoError(t, err)

	done := make(chan error, 1)
	go func() { done <- run.Wait() }()
	select {
	case err := <-done:
		t.Errorf("oy-expose should be still running: %v", err)
	case <-time.After(1 * time.Second):
	}

	resp, err := http.Get(testURL + "/metrics")
	require.NoError(t, err)
	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "node_textfile_scrape_error 0")

	os.Remove(tmpfile.Name())

	resp, err = http.Get(testURL + "/metrics")
	require.NoError(t, err)
	body, err = io.ReadAll(resp.Body)
	require.NoError(t, err)
	require.Contains(t, string(body), "node_textfile_scrape_error 1")

	select {
	case err := <-done:
		t.Errorf("oy-expose should be still running: %v", err)
	case <-time.After(5 * time.Second):
		require.NoError(t, run.Process.Kill())
		<-done
	}
}
