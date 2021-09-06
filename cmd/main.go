package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/UnderTreeTech/waterdrop/pkg/log"

	"github.com/UnderTreeTech/layout/internal/dao"
	"github.com/UnderTreeTech/layout/internal/server/grpc"
	"github.com/UnderTreeTech/layout/internal/server/http"
	"github.com/UnderTreeTech/layout/internal/service"

	"github.com/UnderTreeTech/waterdrop/pkg/stats"

	"github.com/UnderTreeTech/waterdrop/pkg/trace/jaeger"

	"google.golang.org/grpc/resolver"

	"github.com/UnderTreeTech/waterdrop/pkg/registry/etcd"

	"github.com/UnderTreeTech/waterdrop/pkg/conf"
)

// run: go run main.go -conf=../configs/application.toml
func main() {
	flag.Parse()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)

	conf.Init()

	logCfg := &log.Config{}
	if err := conf.Unmarshal("log", logCfg); err != nil {
		panic(fmt.Sprintf("parse log config fail, err msg %s", err.Error()))
	}
	defer log.New(logCfg).Sync()

	// you can commnet this line, then it will use default mock trace
	defer jaeger.Init()()

	etcdCfg := &etcd.Config{}
	if err := conf.Unmarshal("etcd", etcdCfg); err != nil {
		panic(fmt.Sprintf("unmarshal etcd config fail, err msg %s", err.Error()))
	}
	etcd := etcd.New(etcdCfg)
	resolver.Register(etcd)

	dao := dao.New()
	s := service.New(dao)
	http := http.New(s)
	rpc := grpc.New(s)

	etcd.Register(context.Background(), rpc.ServiceInfo)
	etcd.Register(context.Background(), http.ServiceInfo)
	si, err := stats.StartStats()
	if err != nil {
		panic(fmt.Sprintf("start stats fail, err msg is %s", err.Error()))
	}
	etcd.Register(context.Background(), si)

	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	etcd.Close()
	http.Server.Stop(ctx)
	rpc.Server.Stop(ctx)
	s.Close()
}
