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
	"errors"
	"flag"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

type Config struct {
	Brokers string
	Topic   string
	Total   int
}

type consumer struct {
	ready chan struct{}
	done  chan struct{}
	total int
	bytes int
	cost  time.Duration
}

func newConsumer(total int) *consumer {
	return &consumer{
		ready: make(chan struct{}),
		done:  make(chan struct{}),
		total: total,
	}
}

func (c *consumer) Done() {
	<-c.done
}

func (c *consumer) Cost() time.Duration {
	return c.cost
}

func (c *consumer) Bytes() int {
	return c.bytes
}

func (c *consumer) Setup(sarama.ConsumerGroupSession) error {
	close(c.ready)
	return nil
}

func (c *consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (c *consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	var total int
	var start time.Time

	for {
		select {
		case message, ok := <-claim.Messages():
			if !ok {
				return nil
			}
			if start.IsZero() {
				start = time.Now()
			}

			session.MarkMessage(message, "")

			total++
			c.bytes += len(message.Value)
			if total >= c.total {
				close(c.done)
				c.cost = time.Since(start)
				return nil
			}

		case <-session.Context().Done():
			return nil
		}
	}
}

type Client struct {
	ctx    context.Context
	cancel context.CancelFunc

	conf     Config
	cli      sarama.ConsumerGroup
	consumer *consumer
}

func New(conf Config) (*Client, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	group := strconv.Itoa(int(time.Now().Unix()))
	brokers := strings.Split(conf.Brokers, ",")

	cli, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		ctx:      ctx,
		cancel:   cancel,
		conf:     conf,
		cli:      cli,
		consumer: newConsumer(conf.Total),
	}, err
}

func (c *Client) Run() error {
	topics := []string{c.conf.Topic}

	go func() {
		c.consumer.Done()
		c.cancel()
	}()

	for {
		if err := c.cli.Consume(c.ctx, topics, c.consumer); err != nil {
			if errors.Is(err, sarama.ErrClosedConsumerGroup) {
				return nil
			}
			return err
		}
		if c.ctx.Err() != nil {
			return c.ctx.Err()
		}
	}
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func main() {
	var c Config
	flag.StringVar(&c.Brokers, "brokers", "localhost:9092", "kafka bootstrap brokers, as a comma separated list")
	flag.StringVar(&c.Topic, "topic", "benchmark", "kafka topics to be consumed")
	flag.IntVar(&c.Total, "total", 100, "require consume messages count")
	flag.Parse()

	client, err := New(c)
	if err != nil {
		log.Fatalf("failed to create client: %v", err)
	}

	client.Run()
	client.Close()
	log.Printf("consumed %d messages in %s, size=%vMB\n", c.Total, client.consumer.cost, client.consumer.bytes/1024/1024)
}
