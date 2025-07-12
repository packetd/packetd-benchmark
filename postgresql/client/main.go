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
	"sync/atomic"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jedib0t/go-pretty/v6/table"

	"github.com/packetd/packetd-benchmark/common"
)

type Config struct {
	DSN      string
	Workers  int
	Total    int
	SQL      string
	Interval time.Duration
}

type Client struct {
	conf Config
	conn *pgxpool.Pool
}

func New(conf Config) *Client {
	config, err := pgxpool.ParseConfig(conf.DSN)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to parse config: %v\n", err)
		os.Exit(1)
	}

	config.MaxConns = int32(conf.Workers)
	conn, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}

	return &Client{
		conf: conf,
		conn: conn,
	}
}

func (c *Client) Close() error {
	c.conn.Close()
	return nil
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
				log.Printf("[%d/%d] sql (%s)\n", counter, c.conf.Total, c.conf.SQL)
			}
			ch <- struct{}{}
		}
		close(ch)
	}()

	start := time.Now()
	wg := sync.WaitGroup{}

	rr := common.NewResourceRecorder()
	rr.Start()

	var rows atomic.Int64
	for i := 0; i < c.conf.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range ch {
				r, err := c.conn.Query(context.Background(), c.conf.SQL)
				if err != nil {
					log.Fatal(err)
				}
				for r.Next() {
					rows.Add(1)
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

	reqTotal := metrics["postgresql_requests_total"]
	printTable(
		"PostgreSQL",
		c.conf.Total,
		c.conf.Workers,
		fmt.Sprintf("%.3fs", elapsed.Seconds()),
		fmt.Sprintf("%.3f", float64(c.conf.Total)/elapsed.Seconds()),
		c.conf.SQL,
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
		"sql",
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
	flag.StringVar(&c.SQL, "sql", "", "sql statement")
	flag.DurationVar(&c.Interval, "interval", 0, "interval per request")
	flag.Parse()

	client := New(c)
	client.Run()

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}
