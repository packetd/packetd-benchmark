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
	"net"
	"time"

	"google.golang.org/grpc"

	"github.com/packetd/packetd-benchmark/grpc/pb"
)

type server struct {
	pb.UnimplementedBenchmarkServer
}

func (s *server) Size(_ context.Context, in *pb.SizeRequest) (*pb.SizeReply, error) {
	duration, _ := time.ParseDuration(in.GetDuration())
	if duration > 0 {
		time.Sleep(duration)
	}
	return &pb.SizeReply{Message: bytes.Repeat([]byte{'x'}, int(in.GetSize()))}, nil
}

func main() {
	addr := flag.String("addr", ":8085", "grpc server address")
	flag.Parse()

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer()
	pb.RegisterBenchmarkServer(s, &server{})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
