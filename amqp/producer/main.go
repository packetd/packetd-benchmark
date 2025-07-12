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
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Config struct {
	DSN         string
	Queue       string
	Total       int
	MessageSize int
}

type Client struct {
	conf Config
	cli  *amqp.Connection
}

func New(conf Config) (*Client, error) {
	conn, err := amqp.Dial(conf.DSN)
	if err != nil {
		return nil, err
	}
	return &Client{
		conf: conf,
		cli:  conn,
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

	if err = c.initQueue(ch); err != nil {
		return err
	}

	err = ch.Confirm(false)
	if err != nil {
		return err
	}
	confirms := ch.NotifyPublish(make(chan amqp.Confirmation, 100))
	go func() {
		for range confirms {
		}
	}()

	payload := bytes.Repeat([]byte{'x'}, c.conf.MessageSize)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < c.conf.Total; i++ {
			err = ch.Publish(
				"",
				c.conf.Queue,
				false,
				false,
				amqp.Publishing{
					ContentType: "application/octet-stream",
					Body:        payload,
				})

			if err != nil {
				log.Fatal(err)
			}
		}
	}()
	wg.Wait()

	time.Sleep(3 * time.Second)
	return nil
}

func (c *Client) initQueue(ch *amqp.Channel) error {
	_, err := ch.QueueDeclarePassive(c.conf.Queue, true, false, false, false, nil)
	if err == nil {
		log.Printf("delete queue %s \n", c.conf.Queue)
		_, err = ch.QueueDelete(c.conf.Queue, false, false, false)
		if err != nil {
			return err
		}
	}

	log.Printf("create queue %s \n", c.conf.Queue)
	_, err = ch.QueueDeclare(c.conf.Queue, true, false, false, false, nil)
	return err
}

func main() {
	var c Config
	flag.StringVar(&c.DSN, "dsn", "amqp://guest:guest@localhost:5672/", "rabbitmq dsn connected to")
	flag.StringVar(&c.Queue, "queue", "benchmark", "rabbitmq queue name")
	flag.IntVar(&c.Total, "total", 100, "require produce messages count")
	flag.IntVar(&c.MessageSize, "message_size", 10240, "rabbitmq message size in bytes")
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
