package utils

import (
	"strings"
	"unsafe"

	"github.com/UnderTreeTech/layout/internal/i18n"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

const (
	defaultLocale = "zh-cn"
	_letters      = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ+-_"
)

// GetLocaleLng get locale language
func GetLocaleLng(lng string) (locale string) {
	// Multiple types, weighted with the quality value syntax:
	// Accept-Language: fr-CH, fr;q=0.9, en;q=0.8, de;q=0.7, *;q=0.5, en-US,en;q=0.5, en, *
	locale = strings.Split(lng, ";")[0]
	locale = strings.TrimSpace(locale)
	if strings.Contains(locale, ",") {
		elements := strings.Split(locale, ",")
		locale = elements[0]
		if strings.Contains(locale, "-") { // 英文系统一返回en
			if strings.Split(locale, "-")[0] == "en" {
				locale = "en"
			}
		}
	}

	// lng default to en if it doesn't have Accept-Language header or accept any language
	if "" == locale || "*" == locale || "zh" == locale {
		locale = defaultLocale
	}

	locale = strings.ToLower(locale)
	return
}

func TranslateError(ctx *gin.Context, err error) (errmsg string) {
	if errs, ok := err.(validator.ValidationErrors); ok {
		lng := GetLocaleLng(ctx.Request.Header.Get("Accept-Language"))
		trans, ok := i18n.GetTranslator().Uni.GetTranslator(lng)
		if !ok {
			trans, _ = i18n.GetTranslator().Uni.GetTranslator("zh")
		}

		errmsg = errs[0].Translate(trans)
		if lngMap, ok := i18n.GetFieldMapping()[lng]; ok {
			if pathMap, ok := lngMap[ctx.Request.URL.Path]; ok {
				if fieldName, ok := pathMap[strings.ToLower(errs[0].Field())]; ok {
					errmsg = strings.Replace(errmsg, errs[0].Field(), fieldName, 1)
				}
			}
		}
	} else {
		errmsg = err.Error()
	}

	return
}

// StripContentType strip content-type
// application/json;charset=utf-8
func StripContentType(contentType string) string {
	i := strings.Index(contentType, ";")
	if i != -1 {
		contentType = contentType[:i]
	}
	return contentType
}

// StringToBytes converts string to byte slice without a memory allocation.
func StringToBytes(s string) (b []byte) {
	if s == "" {
		return nil
	}
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

// BytesToString converts byte slice to string without a memory allocation.
func BytesToString(b []byte) string {
	if len(b) == 0 {
		return ""
	}
	return unsafe.String(unsafe.SliceData(b), len(b))
}
