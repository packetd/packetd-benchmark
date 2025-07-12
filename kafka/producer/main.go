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
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	batchSize = 100
)

type Config struct {
	Brokers     string
	Topic       string
	Total       int
	MessageSize int
}

type Client struct {
	ctx    context.Context
	cancel context.CancelFunc

	conf Config
	cli  *kafka.Writer
}

func New(conf Config) (*Client, error) {
	brokers := strings.Split(conf.Brokers, ",")
	cli := &kafka.Writer{
		Addr:         kafka.TCP(brokers...),
		Topic:        conf.Topic,
		BatchSize:    batchSize,
		BatchTimeout: 1 * time.Second,
		Balancer:     &kafka.LeastBytes{},
		RequiredAcks: kafka.RequireOne,
		MaxAttempts:  3,
	}

	ctx, cancel := context.WithCancel(context.Background())
	return &Client{
		ctx:    ctx,
		cancel: cancel,
		conf:   conf,
		cli:    cli,
	}, nil
}

func (c *Client) initTopic() error {
	brokers := strings.Split(c.conf.Brokers, ",")
	conn, err := kafka.DialContext(c.ctx, "tcp", brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	addr := fmt.Sprintf("%s:%d", controller.Host, controller.Port)
	controllerConn, err := kafka.DialContext(c.ctx, "tcp", addr)
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	partitions, err := controllerConn.ReadPartitions(c.conf.Topic)
	if err != nil {
		return err
	}

	if len(partitions) > 0 {
		err = controllerConn.DeleteTopics(c.conf.Topic)
		if err != nil {
			return err
		}
		log.Printf("delete topic %s \n", c.conf.Topic)
	}

	log.Printf("create topic %s ...\n", c.conf.Topic)
	err = controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             c.conf.Topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
	if err != nil {
		return err
	}

	time.Sleep(time.Second)
	return nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func (c *Client) Run() error {
	if err := c.initTopic(); err != nil {
		return err
	}

	type Message struct {
		Index   int    `json:"index"`
		Payload string `json:"payload"`
	}

	payload := bytes.Repeat([]byte{'x'}, c.conf.MessageSize)
	batch := make([]kafka.Message, 0, batchSize)
	for i := 0; i < c.conf.Total; i++ {
		msg := Message{
			Index:   i,
			Payload: string(payload),
		}
		msgBytes, _ := json.Marshal(msg)
		batch = append(batch, kafka.Message{
			Key:   []byte(strconv.Itoa(msg.Index)),
			Value: msgBytes,
		})

		if len(batch) >= batchSize {
			if err := c.cli.WriteMessages(c.ctx, batch...); err != nil {
				return err
			}
			batch = make([]kafka.Message, 0, batchSize)
		}
	}

	log.Printf("produced %d messages\n", c.conf.Total)
	return nil
}

func main() {
	var c Config
	flag.StringVar(&c.Brokers, "brokers", "localhost:9092", "kafka bootstrap brokers, as a comma separated list")
	flag.StringVar(&c.Topic, "topic", "benchmark", "kafka topics to be consumed")
	flag.IntVar(&c.Total, "total", 100, "require produce messages count")
	flag.IntVar(&c.MessageSize, "message_size", 10240, "kafka message size in bytes")
	flag.Parse()

	client, err := New(c)
	if err != nil {
		log.Fatal(err)
	}

	client.Run()
	client.Close()
}
