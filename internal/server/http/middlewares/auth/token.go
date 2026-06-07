package auth

import (
	"bytes"
	"io"
	"net/http"
	"strings"

	"github.com/UnderTreeTech/layout/internal/ecode"
	"github.com/UnderTreeTech/layout/internal/server/http/model"
	"github.com/UnderTreeTech/layout/internal/utils/reply"
	"github.com/UnderTreeTech/waterdrop/pkg/database/redis"
	"github.com/UnderTreeTech/waterdrop/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/valyala/bytebufferpool"
)

// Token returns a Gin middleware that validates the authentication token carried in every request.
//
// Parameters:
//   - r: Redis client used to look up and validate the token.
//   - noTokenRoutes: exact URL paths (case-insensitive, no leading slash) that bypass token validation.
//
// Flow:
//  1. Skip validation if the request path matches any entry in noTokenRoutes.
//  2. Use io.TeeReader + bytebufferpool.ByteBuffer to read the request body and simultaneously
//     buffer it, then restore c.Request.Body from the buffer so downstream handlers
//     can read the full body again — without double-allocating the body bytes.
//  3. Bind BaseRequest (token field) from the restored body.
//  4. Validate token existence and correctness via Redis.
//  5. Restore the body one final time before calling c.Next().
func Token(r *redis.Redis, noTokenRoutes []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// ── Step 1: skip paths that do not require token validation ──────────
		path := strings.ToLower(strings.Trim(c.Request.URL.Path, "/"))
		for _, v := range noTokenRoutes {
			if strings.ToLower(v) == path {
				c.Next()
				return
			}
		}

		// ── Step 2: buffer body via TeeReader so it can be read twice ────────
		// bytebufferpool.Get() returns a pooled *ByteBuffer; Put() returns it.
		buf := bytebufferpool.Get()
		defer bytebufferpool.Put(buf)

		// tee writes to buf while c.Request.Body is read.
		// After Bind drains tee, buf.B holds the full body bytes;
		// dup is created from buf.B so downstream handlers can read the body again.
		tee := io.TeeReader(c.Request.Body, buf)
		c.Request.Body = io.NopCloser(tee)

		// ── Step 3: bind BaseRequest ─────────────────────────────────────────
		base := &model.BaseRequest{}
		var err error
		// For GET requests with multipart/form-data (no boundary), ShouldBind would
		// fail; fall back to query-string binding.
		if c.Request.Method == http.MethodGet && c.ContentType() == "multipart/form-data" {
			err = c.ShouldBindQuery(base)
		} else {
			err = c.Bind(base) // Bind drains the tee body into buf automatically
		}
		if err != nil {
			log.Error(c.Request.Context(), "bind base request fail", log.String("error", err.Error()))
			c.JSON(http.StatusOK, reply.Reply(c, nil, ecode.InvalidParam))
			c.Abort()
			return
		}
		// After Bind, buf.B contains the full body bytes; create dup for downstream.
		dup := io.NopCloser(bytes.NewReader(buf.B))

		// ── Step 4: validate token via Redis ─────────────────────────────────
		// TODO: Replace the placeholder below with real token validation logic.
		// Example pattern (Redis session lookup):
		//
		//   cacheKey := fmt.Sprintf("token:%s", base.Token)
		//   cacheToken, err := r.Get(c.Request.Context(), cacheKey)
		//   if err != nil {
		//       c.JSON(http.StatusOK, reply.Reply(c, nil, ecode.InternalError))
		//       c.Abort()
		//       return
		//   }
		//   if cacheToken == "" {
		//       c.JSON(http.StatusOK, reply.Reply(c, nil, ecode.InvalidParam))
		//       c.Abort()
		//       return
		//   }
		_ = r // suppress unused-variable warning until real logic is implemented
		if base.Token == "" {
			log.Warn(c.Request.Context(), "missing token")
			c.JSON(http.StatusOK, reply.Reply(c, nil, ecode.InvalidParam))
			c.Abort()
			return
		}

		// ── Step 5: restore body and pass to the next handler ────────────────
		c.Request.Body = dup
		c.Next()
	}
}
