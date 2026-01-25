package middleware

import (
	"context"
	"strings"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"

	"github.com/bifshteksex/hertz-board/internal/service"
)

// Auth returns JWT authentication middleware
func Auth(jwtService *service.JWTService) app.HandlerFunc {
	return func(c context.Context, ctx *app.RequestContext) {
		authHeader := string(ctx.Request.Header.Peek("Authorization"))
		if authHeader == "" {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"error": "Authorization header required",
			})
			ctx.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"error": "Invalid authorization header format",
			})
			ctx.Abort()
			return
		}

		token := parts[1]

		// Validate token
		claims, err := jwtService.ValidateAccessToken(token)
		if err != nil {
			ctx.JSON(consts.StatusUnauthorized, map[string]interface{}{
				"error": "Invalid or expired token",
			})
			ctx.Abort()
			return
		}

		// Set user information in context
		ctx.Set("user_id", claims.UserID)
		ctx.Set("user_email", claims.Email)

		ctx.Next(c)
	}
}
