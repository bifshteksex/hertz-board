package middleware

import (
	"context"
	"strings"

	"github.com/bifshteksex/hertzboard/internal/config"
	"github.com/cloudwego/hertz/pkg/app"
)

// CORS returns a CORS middleware
func CORS(cfg *config.CORSConfig) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		origin := string(ctx.Request.Header.Peek("Origin"))

		// Check if origin is allowed
		allowedOrigin := ""
		for _, allowed := range cfg.AllowedOrigins {
			if allowed == "*" || allowed == origin {
				allowedOrigin = origin
				break
			}
		}

		if allowedOrigin != "" {
			ctx.Response.Header.Set("Access-Control-Allow-Origin", allowedOrigin)
		}

		if cfg.AllowCredentials {
			ctx.Response.Header.Set("Access-Control-Allow-Credentials", "true")
		}

		// Handle preflight requests
		if string(ctx.Request.Method()) == "OPTIONS" {
			ctx.Response.Header.Set("Access-Control-Allow-Methods", strings.Join(cfg.AllowedMethods, ", "))
			ctx.Response.Header.Set("Access-Control-Allow-Headers", strings.Join(cfg.AllowedHeaders, ", "))
			ctx.Response.Header.Set("Access-Control-Max-Age", string(rune(cfg.MaxAge)))
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next(c)
	}
}
