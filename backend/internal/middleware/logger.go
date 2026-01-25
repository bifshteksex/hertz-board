package middleware

import (
	"context"
	"log"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

// Logger logs HTTP requests
func Logger() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		start := time.Now()
		path := string(ctx.Path())
		method := string(ctx.Method())
		requestID := GetRequestID(ctx)

		ctx.Next(c)

		latency := time.Since(start)
		statusCode := ctx.Response.StatusCode()
		clientIP := ctx.ClientIP()

		log.Printf("[%s] %s %s %d %v %s",
			requestID,
			method,
			path,
			statusCode,
			latency,
			clientIP,
		)
	}
}
