package service

import (
	"context"

	"github.com/UnderTreeTech/layout/api/demo"

	"google.golang.org/protobuf/types/known/emptypb"
)

// SayHello is the gRPC handler for the unary Hello RPC.
// It returns an empty response as a demonstration endpoint.
func (s *Service) SayHello(ctx context.Context, req *demo.HelloReq) (reply *emptypb.Empty, err error) {
	reply = &emptypb.Empty{}
	return reply, nil
}

// SayHelloURL is the gRPC-gateway handler for the Hello HTTP endpoint.
// It constructs a greeting message using the name provided in the request.
func (s *Service) SayHelloURL(ctx context.Context, req *demo.HelloReq) (reply *demo.HelloResp, err error) {
	reply = &demo.HelloResp{Content: "Hello " + req.Name}
	return
}
