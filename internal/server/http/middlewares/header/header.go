package header

import (
	"net/http"

	"github.com/UnderTreeTech/layout/internal/utils/reply"

	"github.com/UnderTreeTech/waterdrop/pkg/log"
	"github.com/UnderTreeTech/waterdrop/pkg/server/http/metadata"
	"github.com/UnderTreeTech/waterdrop/pkg/status"
	"github.com/UnderTreeTech/waterdrop/pkg/utils/xstring"
	"github.com/gin-gonic/gin"
)

// Header returns a Gin middleware that validates the Content-Type header.
// For non-GET requests, it requires the Content-Type to be "application/json";
// otherwise it aborts the request with HTTP 400 Bad Request.
func Header() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		if c.Request.Method != http.MethodGet {
			if "application/json" != xstring.StripContentType(c.Request.Header.Get(metadata.HeaderContentType)) {
				log.Warn(ctx, "invalid content-type", log.String("content-type", c.Request.Header.Get(metadata.HeaderContentType)))
				c.AbortWithStatusJSON(http.StatusBadRequest, reply.Reply(c, nil, status.RequestErr))
				return
			}
		}

		c.Next()
	}
}
