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
	"log"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/packetd/packetd-benchmark/common"
	"github.com/packetd/packetd-benchmark/grpc/pb"
)

type Config struct {
	Addr     string
	Workers  int
	Total    int
	BodySize string
	Duration time.Duration
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
	cli  pb.BenchmarkClient
}

func New(conf Config) *Client {
	conn, err := grpc.NewClient(conf.Addr, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	cli := pb.NewBenchmarkClient(conn)
	return &Client{
		conf: conf,
		cli:  cli,
	}
}

func (c *Client) Run() {
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
			log.Printf("[%d/%d] request hello server, size=%s\n", i+1, c.conf.Total, c.conf.BodySize)
			ch <- struct{}{}
		}
		close(ch)
	}()

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
	log.Printf("Total %d requests take %s, qps=%f, bps=%s\n",
		c.conf.Total,
		elapsed,
		float64(c.conf.Total)/elapsed.Seconds(),
		common.HumanizeBit(float64(c.conf.Total*(c.conf.GetBodySize()))/elapsed.Seconds()),
	)
}

func main() {
	var c Config
	flag.IntVar(&c.Workers, "workers", 1, "concurrency workers")
	flag.IntVar(&c.Total, "total", 1, "requests total")
	flag.StringVar(&c.BodySize, "size", "1KB", "request body size")
	flag.DurationVar(&c.Duration, "duration", time.Second, "duration per request")
	flag.StringVar(&c.Addr, "addr", "localhost:8085", "grpc server address")
	flag.Parse()

	client := New(c)
	client.Run()
}
