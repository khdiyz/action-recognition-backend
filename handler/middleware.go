package handler

import "github.com/gin-gonic/gin"

func corsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		origin := ctx.Request.Header.Get("Origin")
		if origin != "" {
			ctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}
		ctx.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		ctx.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, Access-Control-Request-Method, Access-Control-Request-Headers")
		ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH, HEAD")
		ctx.Writer.Header().Set("Access-Control-Max-Age", "3600")

		if ctx.Request.Method == "OPTIONS" {
			ctx.Writer.Header().Set("Content-Length", "0")
			ctx.AbortWithStatus(204)
			return
		}

		ctx.Next()
	}
}

