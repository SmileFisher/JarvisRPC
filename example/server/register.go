package main

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Register struct {
	EtcdClient *clientv3.Client
}

func NewRegister() *Register {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, //如果是集群，就在后面加所有的节点[]string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	return &Register{EtcdClient: cli}
}

func (r *Register) Register() {
	grant, err := r.EtcdClient.Grant(context.TODO(), 15)
	if err != nil {
		fmt.Println(err.Error())
	}

	var ip = "localhost:8889"

	_, err = r.EtcdClient.Put(context.TODO(), "JARVIS/CLOUD/SERVER1/"+ip, ip, clientv3.WithLease(grant.ID))
	if err != nil {
		fmt.Println(err.Error())
	}
	alive, err := r.EtcdClient.KeepAlive(context.TODO(), grant.ID)
	if err != nil {
		fmt.Println(err.Error())
	}

	for {
		select {
		case a := <-alive:
			fmt.Println("keep alive", a)
		}
	}
}
