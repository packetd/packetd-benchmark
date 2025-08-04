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
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"google.golang.org/grpc"

	"github.com/packetd/packetd-benchmark/common"
	"github.com/packetd/packetd-benchmark/grpc/pb"
)

type Config struct {
	Addr     string
	Workers  int
	Total    int
	BodySize string
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
	conn *grpc.ClientConn
	cli  pb.BenchmarkClient
}

func New(conf Config) *Client {
	conn, err := grpc.NewClient(conf.Addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	cli := pb.NewBenchmarkClient(conn)
	return &Client{
		conf: conf,
		conn: conn,
		cli:  cli,
	}
}

func (c *Client) Run() {
	defer c.conn.Close()

	doRequest := func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := c.cli.Size(ctx, &pb.SizeRequest{Size: int64(c.conf.GetBodySize())})
		if err != nil {
			log.Fatalf("size request error: %v\n", err)
		}
	}

	start := time.Now()
	ch := make(chan struct{}, 1)
	go func() {
		for i := 0; i < c.conf.Total; i++ {
			if c.conf.Interval > 0 {
				time.Sleep(c.conf.Interval)
			}
			if common.ShouldLog(c.conf.Total, i) {
				log.Printf("[%d/%d] request hello server, size=%s\n", i+1, c.conf.Total, c.conf.BodySize)
			}
			ch <- struct{}{}
		}
		close(ch)
	}()

	rr := common.NewResourceRecorder()
	rr.Start()

	var wg sync.WaitGroup
	for i := 0; i < c.conf.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range ch {
				doRequest()
			}
		}()
	}
	wg.Wait()

	elapsed := time.Since(start)
	resource := rr.End()

	time.Sleep(time.Second)
	metrics, err := common.RequestProtocolMetrics()
	if err != nil {
		log.Fatal(err)
	}

	reqTotal := metrics["grpc_requests_total"]
	printTable(
		c.conf.Total,
		c.conf.Workers,
		c.conf.BodySize,
		fmt.Sprintf("%.3fs", elapsed.Seconds()),
		fmt.Sprintf("%.3f", float64(c.conf.Total)/elapsed.Seconds()),
		common.HumanizeBit(float64(c.conf.Total*(c.conf.GetBodySize()))/elapsed.Seconds()),
		int(reqTotal),
		fmt.Sprintf("%.3f%%", reqTotal/float64(c.conf.Total)*100),
		fmt.Sprintf("%.3f", resource.CPU),
		fmt.Sprintf("%.3f", resource.Mem/1024/1024),
	)
}

func printTable(columns ...interface{}) {
	header := []interface{}{
		"request",
		"workers",
		"bodySize",
		"elapsed",
		"qps",
		"bps",
		"proto (request)",
		"proto (percent)",
		"cpu (core)",
		"memory (MB)",
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
	flag.IntVar(&c.Workers, "workers", 1, "concurrency workers")
	flag.IntVar(&c.Total, "total", 1, "requests total")
	flag.StringVar(&c.BodySize, "body_size", "1KB", "request body size")
	flag.DurationVar(&c.Interval, "interval", 0, "interval per request")
	flag.StringVar(&c.Addr, "addr", "localhost:8085", "grpc server address")
	flag.Parse()

	client := New(c)
	client.Run()
}
