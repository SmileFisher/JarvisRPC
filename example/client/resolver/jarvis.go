package resolver

import (
	"context"
	"fmt"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/grpc/resolver"
	"sync"
)

func InitJarvisBuilder() {
	resolver.Register(&JarvisBuilder{
		Dis: NewDiscover(),
	})
}

type JarvisBuilder struct {
	Dis *Discover
}

func (j *JarvisBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	fmt.Println("jarvis builder target", target, " cc:", cc, " opt", opts)
	jarvisResolver := &JarvisResolver{
		Schema: target.URL.Scheme,
		Target: target.URL.Host,
		cc:     cc,
		Dis:    j.Dis,
	}
	fmt.Println("res", jarvisResolver)

	get, err := j.Dis.Get(context.TODO(), "JARVIS/CLOUD/"+target.URL.Host)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("etcd get : ", get)
	for _, kv := range get.Kvs {
		jarvisResolver.Store(string(kv.Key), string(kv.Value))
	}
	jarvisResolver.updateState()

	go jarvisResolver.watcher()
	return jarvisResolver, nil
}

func (j *JarvisBuilder) Scheme() string {
	return "jarvis"
}

type JarvisResolver struct {
	Schema   string
	Target   string
	cc       resolver.ClientConn
	Dis      *Discover
	AddrPool sync.Map
}

func (j *JarvisResolver) Store(addr, val string) {
	j.AddrPool.Store(addr, val)
}

func (j *JarvisResolver) Delete(addr string) {
	j.AddrPool.Delete(addr)
}

func (j *JarvisResolver) ResolveNow(options resolver.ResolveNowOptions) {
	fmt.Println("ResolveNow")
}

func (j *JarvisResolver) watcher() {

	watch := j.Dis.EtcdClient.Watch(context.TODO(), j.Target, clientv3.WithPrefix())
	for {
		select {
		case event := <-watch:
			for _, e := range event.Events {
				if e.Type == clientv3.EventTypePut {
					fmt.Println("put ", string(e.Kv.Key))
					j.Store(string(e.Kv.Key), string(e.Kv.Value))
					j.updateState()
				} else if e.Type == clientv3.EventTypeDelete {
					fmt.Println("delete ", string(e.Kv.Key))
					fmt.Println(string(e.Kv.Value))
					j.Delete(string(e.Kv.Key))
					j.updateState()
				}
			}
		}
	}
}

func (j *JarvisResolver) updateState() {
	var address []resolver.Address
	j.AddrPool.Range(func(key, value any) bool {
		address = append(address, resolver.Address{
			Addr: value.(string),
			//ServerName:         j.Target,
			Attributes:         nil,
			BalancerAttributes: nil,
		})
		fmt.Println("xxxxxxx ", value.(string))
		return true
	})
	fmt.Println("addrss", address)
	j.cc.UpdateState(resolver.State{Addresses: address})
}

func (j *JarvisResolver) Close() {
	fmt.Println("resolver close")
	return
}
