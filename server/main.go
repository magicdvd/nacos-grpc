package main

import (
	"context"
	"flag"
	"fmt"
	"net"

	"github.com/magicdvd/nacos-client"
	"github.com/magicdvd/nacos-grpc/pb"
	"google.golang.org/grpc"
)

var (
	port   = flag.Uint("port", 50001, "port")
	weight = flag.Float64("weight", 10.0, "serviceWeight")
)

type hello struct {
	idx    uint
	weight float64
}

func (c *hello) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloResponse, error) {
	a := new(pb.HelloResponse)
	a.Name = fmt.Sprintf("Hello %s %d, weight: %f", in.Name, c.idx, c.weight)
	a.Port = uint32(c.idx)
	a.Weight = c.weight
	return a, nil
}

func main() {
	flag.Parse()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		fmt.Println(err)
		return
	}
	srv := grpc.NewServer()
	pb.RegisterHelloCenterServer(srv, &hello{
		idx:    *port,
		weight: *weight,
	})
	client, err := nacos.NewServiceClient("http://nacos:nacos@127.0.0.1:8848/nacos", nacos.LogLevel("debug"), nacos.EnableHTTPRequestLog(false))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = client.RegisterInstance("", *port, "hello", nacos.ParamWeight(*weight))
	if err != nil {
		fmt.Println(err)
		return
	}
	err = srv.Serve(l)
	if err != nil {
		fmt.Println(err)
		return
	}
}
