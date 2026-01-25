package middleware

import (
	"context"
	"log"
	"runtime/debug"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// Recovery recovers from panics and returns 500
func Recovery() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(ctx)
				stack := string(debug.Stack())

				log.Printf("[%s] PANIC: %v\n%s", requestID, err, stack)

				ctx.JSON(consts.StatusInternalServerError, map[string]interface{}{
					"error":      "Internal server error",
					"request_id": requestID,
				})
				ctx.Abort()
			}
		}()

		ctx.Next(c)
	}
}
