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

package main

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	addr := flag.String("addr", "localhost:8083", "http server address")
	flag.Parse()

	http.HandleFunc("/benchmark", func(w http.ResponseWriter, r *http.Request) {
		duration, _ := time.ParseDuration(r.FormValue("duration"))
		size, _ := strconv.Atoi(r.FormValue("size"))
		status, _ := strconv.Atoi(r.FormValue("status"))

		log.Printf("request from %s, duration=%v, size=%v, status=%v\n", r.RemoteAddr, duration, size, status)
		if duration > 0 {
			time.Sleep(duration)
		}
		if status >= 200 && status <= 599 {
			w.WriteHeader(status)
			w.Write(bytes.Repeat([]byte{'x'}, size))
			return
		}
		if size > 0 {
			w.Write(bytes.Repeat([]byte{'x'}, size))
		}
	})

	log.Printf("server listening on %s\n", *addr)
	http.ListenAndServe(*addr, nil)
}
