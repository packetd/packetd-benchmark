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

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	DSN        string
	Workers    int
	Total      int
	Database   string
	Collection string
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
			log.Printf("[%d/%d] collection (%s.%s)\n", counter, c.conf.Total, c.conf.Database, c.conf.Collection)
			ch <- struct{}{}
		}
		close(ch)
	}()

	start := time.Now()
	wg := sync.WaitGroup{}

	collection := c.cli.Database(c.conf.Database).Collection(c.conf.Collection)
	for i := 0; i < c.conf.Workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for range ch {
				r, err := collection.Find(c.ctx, bson.D{})
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
	log.Printf("Total %d requests take %s, qps=%f\n",
		c.conf.Total,
		elapsed,
		float64(c.conf.Total)/elapsed.Seconds(),
	)
}

func main() {
	var c Config
	flag.StringVar(&c.DSN, "dsn", "", "mysql server dsn")
	flag.IntVar(&c.Workers, "workers", 1, "concurrency workers")
	flag.IntVar(&c.Total, "total", 1, "requests total")
	flag.StringVar(&c.Database, "database", "", "database name")
	flag.StringVar(&c.Collection, "collection", "", "collection name")
	flag.DurationVar(&c.Interval, "interval", time.Second, "interval between requests")
	flag.Parse()

	client := New(c)
	client.Run()

	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}
