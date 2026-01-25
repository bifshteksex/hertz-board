package middleware

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/google/uuid"
)

const RequestIDHeader = "X-Request-ID"

// RequestID adds a unique request ID to each request
func RequestID() app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		requestID := string(ctx.Request.Header.Peek(RequestIDHeader))
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx.Response.Header.Set(RequestIDHeader, requestID)
		ctx.Set("request_id", requestID)
		ctx.Next(c)
	}
}

// GetRequestID retrieves request ID from context
func GetRequestID(ctx *app.RequestContext) string {
	if requestID, exists := ctx.Get("request_id"); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
