package main

import (
	"context"
	"fmt"

	nacosgrpc "github.com/magicdvd/nacos-grpc"
	"github.com/magicdvd/nacos-grpc/balancer/weightedroundrobin"
	"github.com/magicdvd/nacos-grpc/pb"
	"google.golang.org/grpc"
)

func main() {
	con, err := grpc.Dial(nacosgrpc.Target("http://nacos:nacos@127.0.0.1:8848/nacos", "hello"), grpc.WithInsecure(), grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, weightedroundrobin.Name)))
	if err != nil {
		fmt.Println(err)
		return
	}
	client := pb.NewHelloCenterClient(con)
	a := make(map[string]int)
	for i := 0; i < 999; i++ {
		resp, err := client.SayHello(context.TODO(), &pb.HelloRequest{Name: "me"})
		if err != nil {
			fmt.Println(err)
			return
		}
		k := fmt.Sprintf("%d-%f", resp.Port, resp.Weight)
		if v, ok := a[k]; ok {
			v++
			a[k] = v
		} else {
			a[k] = 1
		}
	}
	for k, v := range a {
		fmt.Println(k, "count:", v)
	}
}
