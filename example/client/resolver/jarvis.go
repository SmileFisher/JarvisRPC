package resolver

import (
	"fmt"
	"google.golang.org/grpc/resolver"
	"time"
)

func InitJarvisBuilder() {
	resolver.Register(&JarvisBuilder{})
}

type JarvisBuilder struct {
}

func (j *JarvisBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {

	fmt.Println("jarvis builder target", target, " cc:", cc, " opt", opts)

	jarvisResolver := &JarvisResolver{
		Schema: target.URL.Scheme,
		Target: target.URL.Host,
		cc:     cc,
	}
	fmt.Println("res", jarvisResolver)
	jarvisResolver.updateState()
	go jarvisResolver.watcher()
	return jarvisResolver, nil
}

func (j *JarvisBuilder) Scheme() string {
	return "jarvis"
}

type JarvisResolver struct {
	Schema string
	Target string
	cc     resolver.ClientConn
}

func (j *JarvisResolver) ResolveNow(options resolver.ResolveNowOptions) {
	j.updateState()
}

func (j *JarvisResolver) watcher() {

	for {
		time.Sleep(3 * time.Second)
		fmt.Println("water .... ")
		j.cc.UpdateState(resolver.State{Addresses: []resolver.Address{
			{
				Addr:               "localhost:8889",
				ServerName:         "SERVER1",
				Attributes:         nil,
				BalancerAttributes: nil,
			},
		}})
	}
}

func (j *JarvisResolver) updateState() {

	fmt.Println("update state .... ")
	j.cc.UpdateState(resolver.State{Addresses: []resolver.Address{
		{
			Addr:               "localhost:8888",
			ServerName:         "SERVER1",
			Attributes:         nil,
			BalancerAttributes: nil,
		},
	}})
}

func (j *JarvisResolver) Close() {
	fmt.Println("resolver close")
	return
}
