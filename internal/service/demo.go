package service

import (
	"context"

	"github.com/UnderTreeTech/layout/api/demo"

	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Service) SayHello(ctx context.Context, req *demo.HelloReq) (reply *emptypb.Empty, err error) {
	reply = &emptypb.Empty{}
	return reply, nil
}
func (s *Service) SayHelloURL(ctx context.Context, req *demo.HelloReq) (reply *demo.HelloResp, err error) {
	reply = &demo.HelloResp{Content: "Hello " + req.Name}
	return
}
