package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"grpc-client-demo/message"
	"grpc-client-demo/resolver"
	"strconv"
	"time"
)

func main() {
	resolver.InitJarvisBuilder()
	for i := 1; i < 100; i++ {
		time.Sleep(5 * time.Second)
		call(i)
	}
}

var dialContext *grpc.ClientConn

func call(i int) {
	var err error
	if dialContext == nil {
		dialContext, err = grpc.DialContext(context.Background(), "jarvis:SERVER"+strconv.Itoa(1),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
			grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`))
		if err != nil {
			fmt.Println(err.Error())
			return
		}
	}

	hello, err := message.NewGreeterClient(dialContext).SayHello(context.TODO(), &message.HelloRequest{
		Name: "hello",
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(hello.Message)
}
