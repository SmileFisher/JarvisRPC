package resolver

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Discover struct {
	EtcdClient *clientv3.Client
}

func NewDiscover() *Discover {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, //如果是集群，就在后面加所有的节点[]string{"localhost:2379", "localhost:22379", "localhost:32379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		fmt.Println(err.Error())
	}
	return &Discover{EtcdClient: cli}
}

//服务存储方式 JARVIS/CLOUD/DEMOA/{HOST:PORT}:{HOST:PORT}

func (d *Discover) Get(ctx context.Context, prefix string) (*clientv3.GetResponse, error) {
	slist, err := d.EtcdClient.Get(ctx, prefix, clientv3.WithPrefix())

	if err != nil {
		fmt.Println(err.Error())
	}
	return slist, nil
}
