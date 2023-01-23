/*
Copyright 2023 Adam B Kaplan

SPDX-License-Id: Apache-2.0
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/adambkaplan/go-k8s-middleware/go-grpc-k8s-auth/authn"
	pb "github.com/adambkaplan/go-k8s-middleware/go-grpc-k8s-auth/example/exampleproto"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: req.GetName()}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.StreamInterceptor(authn.TokenReviewStreamServerInterceptor()),
		grpc.UnaryInterceptor(authn.TokenReviewUnaryServerInterceptor()))

	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
