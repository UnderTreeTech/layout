package ecode

import "github.com/UnderTreeTech/waterdrop/pkg/status"

var (
	InternalError = status.New(10000000, "网络开小差，请稍后重试")
	InvalidParam  = status.New(10000001, "参数错误")
)

var enRequestErrMsg = errors{
	InternalError.Code(): "Network is unavailable, please have a try later",
	InvalidParam.Code():  "Invalid parameters, please check your request",
}

var zhRequestErrMsg = errors{
	InternalError.Code(): InternalError.Message(),
	InvalidParam.Code():  InvalidParam.Message(),
}
