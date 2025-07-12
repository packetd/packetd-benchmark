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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/packetd/packetd-benchmark/common"
)

type Config struct {
	DSN        string
	Workers    int
	Total      int
	Database   string
	Collection string
	Limit      int64
	Interval   time.Duration
}

type Client struct {
	ctx    context.Context
	cancel context.CancelFunc
	conf   Config
	cli    *mongo.Client
}

func New(conf Config) *Client {
	cli, err := mongo.Connect(context.Background(), options.Client().ApplyURI(conf.DSN))
	if err != nil {
		log.Fatal(err)
	}

	options.Client().SetMaxConnecting(64)
	options.Client().SetMinPoolSize(64)

	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		ctx:    ctx,
		cancel: cancel,
		conf:   conf,
		cli:    cli,
	}
}

func (c *Client) Close() error {
	return c.cli.Disconnect(context.Background())
}

func (c *Client) Run() {
	ch := make(chan struct{}, 1)
	go func() {
		var counter int
		for i := 0; i < c.conf.Total; i++ {
			counter++
			if c.conf.Interval > 0 {
				time.Sleep(c.conf.Interval)
			}
			if common.ShouldLog(c.conf.Total, i) {
				log.Printf("[%d/%d] collection (%s.%s)\n", counter, c.conf.Total, c.conf.Database, c.conf.Collection)
			}
			ch <- struct{}{}
		}
		close(ch)
	}()

	start := time.Now()
	wg := sync.WaitGroup{}

	rr := common.NewResourceRecorder()
	rr.Start()

	collection := c.cli.Database(c.conf.Database).Collection(c.conf.Collection)
	for i := 0; i < c.conf.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range ch {
				opt := options.Find()
				opt.Limit = &c.conf.Limit
				r, err := collection.Find(c.ctx, bson.D{}, opt)
				if err != nil {
					log.Fatal(err)
				}

				for r.Next(c.ctx) {
				}
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

	reqTotal := metrics["mongodb_requests_total"]
	printTable(
		"MongoDB",
		c.conf.Total,
		c.conf.Workers,
		fmt.Sprintf("%.3fs", elapsed.Seconds()),
		fmt.Sprintf("%.3f", float64(c.conf.Total)/elapsed.Seconds()),
		c.conf.Limit,
		int(reqTotal),
		fmt.Sprintf("%.3f%%", reqTotal/float64(c.conf.Total)*100),
		fmt.Sprintf("%.3f", resource.CPU),
		fmt.Sprintf("%.3f", resource.Mem/1024/1024),
	)
}

func printTable(columns ...interface{}) {
	header := []interface{}{
		"proto",
		"request",
		"workers",
		"elapsed",
		"qps",
		"limit",
		"proto/request",
		"proto/percent",
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
	flag.StringVar(&c.DSN, "dsn", "", "mysql server dsn")
	flag.IntVar(&c.Workers, "workers", 1, "concurrency workers")
	flag.IntVar(&c.Total, "total", 1, "requests total")
	flag.StringVar(&c.Database, "database", "", "database name")
	flag.StringVar(&c.Collection, "collection", "", "collection name")
	flag.DurationVar(&c.Interval, "interval", 0, "interval per request")
	flag.Int64Var(&c.Limit, "limit", 0, "records count")
	flag.Parse()

	client := New(c)
	client.Run()

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}
