package grpc

import (
	"fmt"
	"net"

	"github.com/UnderTreeTech/layout/api/demo"
	"github.com/UnderTreeTech/layout/internal/service"

	"github.com/UnderTreeTech/waterdrop/pkg/server/rpc/config"

	"github.com/UnderTreeTech/waterdrop/pkg/server/rpc/server"

	"github.com/UnderTreeTech/waterdrop/pkg/server/rpc/interceptors"

	"github.com/UnderTreeTech/waterdrop/pkg/utils/xnet"

	"github.com/UnderTreeTech/waterdrop/pkg/conf"

	"google.golang.org/grpc"

	"github.com/UnderTreeTech/waterdrop/pkg/registry"
)

// ServerInfo holds the gRPC server instance and its service registration metadata.
type ServerInfo struct {
	Server      *server.Server
	ServiceInfo *registry.ServiceInfo
}

// New creates, configures, and starts the gRPC server.
// It loads server configuration, registers all gRPC service implementations,
// attaches the validation interceptor, and returns the running ServerInfo.
func New(s *service.Service) *ServerInfo {
	srvConfig := &config.ServerConfig{}
	parseConfig("server.rpc", srvConfig)
	if srvConfig.WatchConfig {
		conf.OnChange(func(config *conf.Config) {
			parseConfig("server.rpc", srvConfig)
		})
	}

	server := server.New(srvConfig)
	registerServers(server.Server(), s)

	server.Use(interceptors.ValidateForUnaryServer())

	addr := server.Start()
	_, port, _ := net.SplitHostPort(addr.String())
	serviceInfo := &registry.ServiceInfo{
		Name:    "service.demo.v1",
		Scheme:  "grpc",
		Addr:    fmt.Sprintf("%s://%s:%s", "grpc", xnet.InternalIP(), port),
		Version: "1.0.0",
	}

	return &ServerInfo{Server: server, ServiceInfo: serviceInfo}
}

// registerServers registers all gRPC service implementations onto the given grpc.Server.
func registerServers(g *grpc.Server, s *service.Service) {
	demo.RegisterDemoServer(g, s)
}

// parseConfig unmarshals the named configuration section into srvConfig.
// Panics if the configuration cannot be parsed.
func parseConfig(configName string, srvConfig *config.ServerConfig) {
	if err := conf.Unmarshal(configName, srvConfig); err != nil {
		panic(fmt.Sprintf("unmarshal grpc server config fail, err msg %s", err.Error()))
	}
}
