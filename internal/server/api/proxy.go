package api

import (
	"github.com/acgn-org/onest/internal/source"
	"github.com/gin-gonic/gin"
)

func RealSearchProxy() gin.HandlerFunc {
	proxy := source.RealSearch.NewProxy()
	return func(ctx *gin.Context) {
		ctx.Request.URL.Path = ctx.Param("path")
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
