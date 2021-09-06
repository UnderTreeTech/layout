package http

import (
	"github.com/UnderTreeTech/layout/api/demo"

	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/UnderTreeTech/waterdrop/pkg/log"
	"github.com/gin-gonic/gin"
)

func getUserInfo(ctx *gin.Context) {
	id := ctx.Params.ByName("id")
	log.Info(ctx.Request.Context(), "", log.String("id", id))
	reply, err := svc.SayHelloURL(ctx.Request.Context(), &demo.HelloReq{Name: id})
	if err != nil {
		ctx.JSON(0, Response{Code: 0, Message: "fail", Data: &emptypb.Empty{}})
		return
	}
	ctx.JSON(0, Response{Code: 0, Message: "ok", Data: reply})
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
