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
	"flag"
	"fmt"
	"log"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	DSN   string
	Queue string
	Total int
}

type Client struct {
	conf Config
	cli  *amqp.Connection
}

func New(conf Config) (*Client, error) {
	cli, err := amqp.Dial(conf.DSN)
	if err != nil {
		return nil, err
	}
	return &Client{
		conf: conf,
		cli:  cli,
	}, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func (c *Client) Run() error {
	ch, err := c.cli.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(c.conf.Queue, true, false, false, false, nil)
	if err != nil {
		return err
	}

	err = ch.Qos(100, 0, false)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	var consumer <-chan amqp.Delivery

	id := fmt.Sprintf("consumer-%d", time.Now().Unix())
	consumer, err = ch.Consume(c.conf.Queue, id, false, false, false, false, nil)

	var bytes int
	var n int
	var start time.Time

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case d, ok := <-consumer:
				if !ok {
					return
				}
				d.Ack(false)

				if start.IsZero() {
					start = time.Now()
				}
				n++
				bytes += len(d.Body)
				if n >= c.conf.Total {
					return
				}
			}
		}
	}()

	wg.Wait()
	log.Printf("consumed %d messages in %s, size=%.2fMB\n", c.conf.Total, time.Since(start), float64(bytes)/1024/1024)
	return nil
}

func main() {
	var c Config
	flag.StringVar(&c.DSN, "dsn", "amqp://guest:guest@localhost:5672/", "kafka bootstrap brokers, as a comma separated list")
	flag.StringVar(&c.Queue, "queue", "benchmark", "kafka topics to be consumed")
	flag.IntVar(&c.Total, "total", 100, "require produce messages count")
	flag.Parse()

	client, err := New(c)
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Run(); err != nil {
		log.Fatal(err)
	}
	if err := client.Close(); err != nil {
		log.Fatal(err)
	}
}
