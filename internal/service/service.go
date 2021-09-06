package service

import (
	"fmt"

	"github.com/UnderTreeTech/layout/api/demo"
	"github.com/UnderTreeTech/layout/internal/dao"

	"github.com/UnderTreeTech/waterdrop/pkg/conf"
	"github.com/UnderTreeTech/waterdrop/pkg/server/http/client"
	"github.com/UnderTreeTech/waterdrop/pkg/server/http/config"
	rpcClient "github.com/UnderTreeTech/waterdrop/pkg/server/rpc/client"
	rpcConfig "github.com/UnderTreeTech/waterdrop/pkg/server/rpc/config"
)

type Service struct {
	dao  dao.Dao
	demo demo.DemoClient
	http *client.Client
}

func New(d dao.Dao) *Service {
	cliConf := &rpcConfig.ClientConfig{}
	if err := conf.Unmarshal("client.rpc.demo", cliConf); err != nil {
		panic(fmt.Sprintf("unmarshal demo client config fail, err msg %s", err.Error()))
	}
	rpc := demo.NewDemoClient(rpcClient.New(cliConf).GetConn())

	httpCliConf := &config.ClientConfig{}
	if err := conf.Unmarshal("client.http.app", httpCliConf); err != nil {
		panic(fmt.Sprintf("unmarshal http client config fail, err msg %s", err.Error()))
	}
	httpCli := client.New(httpCliConf)

	return &Service{
		dao:  d,
		http: httpCli,
		demo: rpc,
	}
}

func (s *Service) Close() {
	s.dao.Close()
}
