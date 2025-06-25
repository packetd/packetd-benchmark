// Copyright 2025 The packetd Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func RequestProtocolMetrics(prefix string) (float64, error) {
	rsp, err := http.Get("http://localhost:9091/protocol/metrics")
	if err != nil {
		return 0, err
	}
	defer rsp.Body.Close()

	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		return 0, err
	}

	for _, line := range strings.Split(string(b), "\n") {
		if !strings.HasPrefix(line, prefix) {
			continue
		}

		val := strings.Split(line, " ")[1]
		return strconv.ParseFloat(val, 64)
	}
	return 0, nil
}

func RequestReset() error {
	rsp, err := http.Post("http://localhost:9091/-/reset", "", &bytes.Buffer{})
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	return nil
}
