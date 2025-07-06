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
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func RequestProtocolMetrics() (map[string]float64, error) {
	return doRequest("http://localhost:9091/protocol/metrics")
}

func RequestMetrics() (map[string]float64, error) {
	return doRequest("http://localhost:9091/metrics")
}

func doRequest(url string) (map[string]float64, error) {
	rsp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer rsp.Body.Close()

	b, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, err
	}

	metrics := make(map[string]float64)
	lines := strings.Split(string(b), "\n")
	for _, line := range lines {
		if len(line) == 0 || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.Split(line, " ")
		f, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			continue
		}

		name := strings.Split(parts[0], "{")[0]
		metrics[name] = f
	}
	return metrics, nil
}

type Resource struct {
	CPU float64
	Mem float64
}

type ResourceRecorder struct {
	t     time.Time
	start Resource
}

func NewResourceRecorder() *ResourceRecorder {
	return &ResourceRecorder{}
}

func (r *ResourceRecorder) Start() {
	r.t = time.Now()
	metrics, err := RequestMetrics()
	if err != nil {
		return
	}

	r.start = Resource{
		CPU: metrics["process_cpu_seconds_total"],
	}
}

func (r *ResourceRecorder) End() Resource {
	var resource Resource
	metrics, err := RequestMetrics()
	if err != nil {
		return resource
	}

	cpu := (metrics["process_cpu_seconds_total"] - r.start.CPU) / (time.Now().Sub(r.t).Seconds())
	return Resource{
		CPU: cpu,
		Mem: metrics["process_resident_memory_bytes"],
	}
}
