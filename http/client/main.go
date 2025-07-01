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
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/packetd/packetd-benchmark/common"
)

type Config struct {
	Addr     string
	Workers  int
	Total    int
	BodySize string
	Status   string
	Interval time.Duration
}

func (c Config) GetBodySize() int {
	i, err := common.ParseBytes(c.BodySize)
	if err != nil {
		panic(err)
	}
	return i
}

type Client struct {
	conf Config
	cli  *http.Client
}

func New(conf Config) *Client {
	cli := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout:     time.Minute,
		},
	}
	return &Client{
		conf: conf,
		cli:  cli,
	}
}

func (c *Client) Run() {
	start := time.Now()
	urls := make(chan string, 1)
	statusList := strings.Split(c.conf.Status, ",")

	go func() {
		for i := 0; i < c.conf.Total; i++ {
			u := fmt.Sprintf("http://%s/benchmark?duration=%v&size=%v&status=%v",
				c.conf.Addr,
				c.conf.Interval.String(),
				c.conf.BodySize,
				statusList[i%len(statusList)],
			)

			if (i+1)%(c.conf.Total/10) == 0 || i == c.conf.Total-1 {
				log.Printf("[%d/%d] %s\n", i+1, c.conf.Total, u)
			}
			urls <- u
		}
		close(urls)
	}()

	doRequest := func(u string) error {
		r, _ := http.NewRequest(http.MethodGet, u, &bytes.Buffer{})
		rsp, err := c.cli.Do(r)
		if err != nil {
			return err
		}
		defer rsp.Body.Close()
		io.Copy(io.Discard, rsp.Body)
		return nil
	}

	var wg sync.WaitGroup
	for i := 0; i < c.conf.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for u := range urls {
				if err := doRequest(u); err != nil {
					log.Fatal(err)
				}
			}
		}()
	}
	wg.Wait()

	time.Sleep(time.Second)
	reqTotal, _ := common.RequestProtocolMetrics("http_requests_total")
	elapsed := time.Since(start)
	printTable(
		"HTTP",
		c.conf.Total,
		c.conf.Workers,
		c.conf.BodySize,
		fmt.Sprintf("%.3fs", elapsed.Seconds()),
		fmt.Sprintf("%.3f", float64(c.conf.Total)/elapsed.Seconds()),
		common.HumanizeBit(float64(c.conf.Total*(c.conf.GetBodySize()))/elapsed.Seconds()),
		reqTotal,
		fmt.Sprintf("%.3f%%", reqTotal/float64(c.conf.Total)*100),
	)
	_ = common.RequestReset()
}

func printTable(columns ...interface{}) {
	header := []interface{}{
		"Proto",
		"Request",
		"Workers",
		"BodySize",
		"Elapsed",
		"QPS",
		"bps",
		"Proto/Metrics",
		"Proto/Percent",
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(header)
	t.AppendRow(columns)
	t.AppendSeparator()
	t.Render()
}

func main() {
	var c Config
	flag.StringVar(&c.Addr, "addr", "localhost:8083", "http server address")
	flag.IntVar(&c.Workers, "workers", 1, "concurrency workers")
	flag.IntVar(&c.Total, "total", 1, "requests total")
	flag.StringVar(&c.BodySize, "body_size", "1KB", "request body size")
	flag.DurationVar(&c.Interval, "interval", 0, "interval per request")
	flag.StringVar(&c.Status, "status", "200", "http response status")
	flag.Parse()

	client := New(c)
	client.Run()
}
