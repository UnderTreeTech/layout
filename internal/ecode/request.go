package ecode

import "github.com/UnderTreeTech/waterdrop/pkg/status"

// Pre-defined business error codes.
var (
	// InternalError indicates an internal server error or network issue.
	InternalError = status.New(100000, "网络开小差，请稍后重试")
	// InvalidParam indicates the request parameters are invalid or failed validation.
	InvalidParam = status.New(100001, "参数错误")
)

// enRequestErrMsg holds the English translations for each error code.
var enRequestErrMsg = errors{
	InternalError.Code(): "Network is unavailable, please have a try later",
	InvalidParam.Code():  "Invalid parameters, please check your request",
}

// zhRequestErrMsg holds the Chinese translations for each error code.
var zhRequestErrMsg = errors{
	InternalError.Code(): InternalError.Message(),
	InvalidParam.Code():  InvalidParam.Message(),
}
