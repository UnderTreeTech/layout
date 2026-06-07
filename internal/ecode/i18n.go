package ecode

import (
	"github.com/UnderTreeTech/layout/internal/utils"

	"github.com/gin-gonic/gin"
)

// errors is a mapping of error code to its human-readable message string.
type errors map[int]string

// ecodes stores the registered error messages grouped by locale (e.g. "en", "zh-cn").
var ecodes = make(map[string]errors)

// Register register error msg
func Register(codes map[string]errors) {
	for locale, msgs := range codes {
		if existLocaleCodes, ok := ecodes[locale]; ok {
			for code, msg := range msgs {
				existLocaleCodes[code] = msg
			}
			ecodes[locale] = existLocaleCodes
		} else {
			ecodes[locale] = msgs
		}
	}
}

// Message return the errmsg associate with the errcode.
// When code is 0 (success), it always returns "ok" regardless of locale.
// If the locale is not registered, it falls back to the Chinese ("zh-cn") locale.
func Message(ctx *gin.Context, code int) string {
	if code == 0 {
		return "ok"
	}

	locale := utils.GetLocaleLng(ctx.Request.Header.Get("Accept-Language"))

	if localCodes, ok := ecodes[locale]; ok {
		return localCodes[code]
	}

	// fallback to zh-cn if the requested locale is not registered
	if zhCodes, ok := ecodes["zh-cn"]; ok {
		return zhCodes[code]
	}

	return ""
}

func init() {
	// 注册自定义错误码，相同错误码会覆盖全局，达到更换文案目的
	Register(map[string]errors{
		"en":    enRequestErrMsg,
		"zh-cn": zhRequestErrMsg,
	})
}
