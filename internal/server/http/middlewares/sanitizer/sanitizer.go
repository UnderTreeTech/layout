package sanitizer

import (
	"reflect"
	"unicode"

	"github.com/UnderTreeTech/waterdrop/pkg/utils/xslice"
	"github.com/gin-gonic/gin"

	"github.com/microcosm-cc/bluemonday"
)

// Sanitizer wraps a bluemonday.Policy and provides recursive XSS sanitization
// for arbitrary request structs via reflection.
type Sanitizer struct {
	p *bluemonday.Policy
}

// skipUrls is a list of URL paths that bypass sanitization entirely.
// Add paths here when the endpoint intentionally accepts raw HTML content.
var skipUrls = []string{}

// NewSanitizer creates a new Sanitizer using the bluemonday UGC policy.
// The UGC policy allows a safe subset of HTML typically used in user-generated content.
func NewSanitizer() *Sanitizer {
	sanitizer := &Sanitizer{
		p: bluemonday.UGCPolicy(),
	}

	// todo: customize allow elements
	// sanitize.p.AllowElements("style")

	return sanitizer
}

// Sanitize performs XSS sanitization on all string fields of req.
// Requests whose URL path is listed in skipUrls are skipped entirely.
func (s *Sanitizer) Sanitize(ctx *gin.Context, req interface{}) {
	if xslice.ContainString(skipUrls, ctx.Request.URL.Path) {
		return
	}
	s.sanitize(ctx, reflect.ValueOf(req))
}

// sanitize recursively traverses v and sanitizes all string values against XSS attacks.
// It handles Ptr/Interface (with nil guard), Struct, Slice/Array, Map, and String kinds.
func (s *Sanitizer) sanitize(ctx *gin.Context, v reflect.Value) {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		// guard against nil pointer/interface to prevent panic on v.Elem()
		if v.IsNil() {
			return
		}
		s.sanitize(ctx, v.Elem())
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			s.sanitize(ctx, v.Field(i))
		}
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			s.sanitize(ctx, v.Index(i))
		}
	case reflect.Map:
		// sanitize string values inside map fields to prevent XSS bypass via map fields.
		// For string values: sanitize in-place via SetMapIndex.
		// For Ptr/Struct values: create a temporary addressable copy, recursively sanitize it,
		// and write back — so nested struct fields inside map values are also cleaned.
		for _, key := range v.MapKeys() {
			val := v.MapIndex(key)
			// dereference interface wrapper (map values are wrapped in interface{})
			unwrapped := val
			if unwrapped.Kind() == reflect.Interface {
				unwrapped = unwrapped.Elem()
			}
			switch unwrapped.Kind() {
			case reflect.String:
				str := unwrapped.String()
				if len(str) == 0 || isNumeric(str) {
					continue
				}
				v.SetMapIndex(key, reflect.ValueOf(s.p.Sanitize(str)))
			case reflect.Ptr:
				// map values are not addressable; recursive call works because
				// we dereference the pointer to get an addressable struct
				if !unwrapped.IsNil() {
					s.sanitize(ctx, unwrapped)
				}
			case reflect.Struct:
				// map values are not settable, so make a new pointer copy,
				// sanitize it, and write back
				ptr := reflect.New(unwrapped.Type())
				ptr.Elem().Set(unwrapped)
				s.sanitize(ctx, ptr.Elem())
				v.SetMapIndex(key, ptr.Elem())
			}
		}
	case reflect.String:
		// html.EscapeString is lighter, bluemonday is more comprehensive
		// Both sanitize will replace &
		// v.SetString(html.EscapeString(v.String()))
		if !v.CanSet() {
			return
		}

		if isNumeric(v.String()) || len(v.String()) == 0 {
			return
		}

		v.SetString(s.p.Sanitize(v.String()))
	}
}

// isNumeric Checks if the string contains only digits. A decimal point is not a digit and returns false.
func isNumeric(s string) bool {
	if s == "" {
		return false
	}
	for _, v := range s {
		if !unicode.IsDigit(v) {
			return false
		}
	}
	return true
}
