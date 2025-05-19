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

	config.MinIdleConns = 10
	config.MaxConns = 50

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
			log.Printf("[%d/%d] sql (%s)\n", counter, c.conf.Total, c.conf.SQL)
			ch <- struct{}{}
		}
		close(ch)
	}()

	start := time.Now()
	wg := sync.WaitGroup{}

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
	log.Printf("Total %d requests take %s, qps=%f, rows=%d\n",
		c.conf.Total,
		elapsed,
		float64(c.conf.Total)/elapsed.Seconds(),
		rows.Load(),
	)
}

func main() {
	var c Config
	flag.StringVar(&c.DSN, "dsn", "", "mysql server dsn")
	flag.IntVar(&c.Workers, "workers", 1, "concurrency workers")
	flag.IntVar(&c.Total, "total", 1, "requests total")
	flag.StringVar(&c.SQL, "sql", "", "sql statement")
	flag.DurationVar(&c.Interval, "interval", time.Second, "interval between requests")
	flag.Parse()

	client := New(c)
	client.Run()

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}
