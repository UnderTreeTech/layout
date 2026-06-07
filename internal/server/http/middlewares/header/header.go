package header

import (
	"net/http"

	"github.com/UnderTreeTech/layout/internal/ecode"
	"github.com/UnderTreeTech/layout/internal/utils/reply"

	"github.com/UnderTreeTech/waterdrop/pkg/log"
	"github.com/UnderTreeTech/waterdrop/pkg/server/http/metadata"
	"github.com/UnderTreeTech/waterdrop/pkg/utils/xstring"
	"github.com/gin-gonic/gin"
)

// allowedContentTypes lists the Content-Type values accepted for request body methods.
// - application/json: standard JSON API requests
// - multipart/form-data: file upload requests
// - application/x-www-form-urlencoded: HTML form submissions
var allowedContentTypes = map[string]struct{}{
	"application/json":                  {},
	"multipart/form-data":               {},
	"application/x-www-form-urlencoded": {},
}

// Header returns a Gin middleware that validates the Content-Type header.
// Only POST, PUT, and PATCH requests require a valid Content-Type;
// other methods (GET, HEAD, DELETE, OPTIONS, etc.) are allowed through without check.
// This avoids blocking CORS preflight OPTIONS requests and body-less DELETE/HEAD requests.
// Accepted Content-Types: application/json, multipart/form-data, application/x-www-form-urlencoded.
// On invalid Content-Type, the request is aborted with HTTP 200 + ecode.InvalidParam
// to be consistent with the rest of the project's error response format.
func Header() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()

		// Only enforce Content-Type on methods that carry a request body.
		// Original code used "!= GET" which also blocked OPTIONS (CORS preflight),
		// HEAD, and body-less DELETE requests.
		switch c.Request.Method {
		case http.MethodPost, http.MethodPut, http.MethodPatch:
			ct := xstring.StripContentType(c.Request.Header.Get(metadata.HeaderContentType))
			if _, ok := allowedContentTypes[ct]; !ok {
				log.Warn(ctx, "invalid content-type", log.String("content-type", c.Request.Header.Get(metadata.HeaderContentType)))
				c.AbortWithStatusJSON(http.StatusOK, reply.Reply(c, nil, ecode.InvalidParam))
				return
			}
		}

		c.Next()
	}
}
