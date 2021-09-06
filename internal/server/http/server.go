package http

import (
	"fmt"
	"net"

	"github.com/UnderTreeTech/layout/internal/service"

	"github.com/UnderTreeTech/waterdrop/pkg/server/http/config"

	"github.com/UnderTreeTech/waterdrop/pkg/server/http/server"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"

	"github.com/UnderTreeTech/waterdrop/pkg/conf"

	"github.com/UnderTreeTech/waterdrop/pkg/utils/xnet"

	"github.com/UnderTreeTech/waterdrop/pkg/registry"
)

type ServerInfo struct {
	Server      *server.Server
	ServiceInfo *registry.ServiceInfo
}

var svc *service.Service

func New(s *service.Service) *ServerInfo {
	srvConfig := &config.ServerConfig{}
	parseConfig("server.http", srvConfig)
	if srvConfig.WatchConfig {
		conf.OnChange(func(config *conf.Config) {
			parseConfig("server.http", srvConfig)
		})
	}

	server := server.New(srvConfig)
	registerMiddlewares(server)
	router(server)
	svc = s
	addr := server.Start()
	_, port, _ := net.SplitHostPort(addr.String())
	serviceInfo := &registry.ServiceInfo{
		Name:    "server.http.example",
		Scheme:  "http",
		Addr:    fmt.Sprintf("%s://%s:%s", "http", xnet.InternalIP(), port),
		Version: "1.0.0",
	}
	binding.Validator.Engine().(*validator.Validate).SetTagName("validate")
	return &ServerInfo{Server: server, ServiceInfo: serviceInfo}
}

func parseConfig(configName string, srvConfig *config.ServerConfig) {
	if err := conf.Unmarshal(configName, srvConfig); err != nil {
		panic(fmt.Sprintf("unmarshal http server config fail, err msg %s", err.Error()))
	}
}

func registerMiddlewares(s *server.Server) {
	//	register middleware here
}

func router(s *server.Server) {
	g := s.Group("/api")
	{
		g.GET("/app/user/:id", getUserInfo)
	}
}
