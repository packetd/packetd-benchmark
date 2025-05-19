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
	"context"
	"flag"
	"log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/packetd/packetd-benchmark/common"
)

type Config struct {
	Addr     string
	Workers  int
	Total    int
	PoolSize int
	BodySize string
	Cmd      string
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
	cli  *redis.Client
}

func New(conf Config) *Client {
	cli := redis.NewClient(&redis.Options{
		Addr:         conf.Addr,
		DialTimeout:  time.Second,
		ReadTimeout:  time.Second,
		WriteTimeout: time.Second,
		PoolSize:     conf.PoolSize,
	})
	return &Client{
		conf: conf,
		cli:  cli,
	}
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func (c *Client) Run() {
	ch := make(chan struct{}, 1)
	go func() {
		var counter int
		for i := 0; i < c.conf.Total; i++ {
			counter++
			log.Printf("[%d/%d] command %s, size=%s\n", counter, c.conf.Total, c.conf.Cmd, c.conf.BodySize)
			ch <- struct{}{}
		}
		close(ch)
	}()

	start := time.Now()
	wg := sync.WaitGroup{}
	for i := 0; i < c.conf.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range ch {
				var err error
				switch c.conf.Cmd {
				case "ping":
					err = c.cmdPing()
				case "set":
					err = c.cmdSet()
				case "get":
					err = c.cmdGet()
				}
				if err != nil {
					log.Fatal(err)
				}
			}
		}()
	}
	wg.Wait()

	elapsed := time.Since(start)
	log.Printf("Total %d requests take %s, qps=%f, bps=%s\n",
		c.conf.Total,
		elapsed,
		float64(c.conf.Total)/elapsed.Seconds(),
		common.HumanizeBit(float64(c.conf.Total)*float64(c.conf.GetBodySize())/elapsed.Seconds()),
	)
}

func (c *Client) cmdPing() error {
	return c.cli.Ping(context.Background()).Err()
}

func (c *Client) cmdSet() error {
	return c.cli.Set(context.Background(), "hello", bytes.Repeat([]byte{'x'}, c.conf.GetBodySize()), 0).Err()
}

func (c *Client) cmdGet() error {
	return c.cli.Get(context.Background(), "hello").Err()
}

func main() {
	var c Config
	flag.StringVar(&c.Addr, "addr", "localhost:6379", "redis server address")
	flag.IntVar(&c.Workers, "workers", 1, "concurrency workers")
	flag.IntVar(&c.Total, "total", 1, "requests total")
	flag.IntVar(&c.PoolSize, "pool_size", 10, "connection pool size")
	flag.StringVar(&c.BodySize, "body_size", "1KB", "request body size")
	flag.StringVar(&c.Cmd, "cmd", "ping", "redis command, options: ping/set/get")
	flag.DurationVar(&c.Interval, "interval", time.Second, "interval between requests")
	flag.Parse()

	client := New(c)
	client.Run()

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}
