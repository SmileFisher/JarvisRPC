package main

import (
	context "context"
	"fmt"
	"google.golang.org/grpc"
	"grpc-service-demo/message"
	"net"
)

type MyServer struct {
	message.UnimplementedGreeterServer
}

func (m *MyServer) SayHello(ctx context.Context, request *message.HelloRequest) (*message.HelloReply, error) {
	fmt.Println(request.Name)

	return &message.HelloReply{
		Message: "world",
	}, nil
}

func main() {

	r := NewRegister()
	go r.Register()

	listen, err := net.Listen("tcp", ":8889")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	server := grpc.NewServer()
	message.RegisterGreeterServer(server, &MyServer{})
	err = server.Serve(listen)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

}
