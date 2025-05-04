package api

import (
	"github.com/acgn-org/onest/realsearch"
	"github.com/gin-gonic/gin"
)

func RealSearchProxy(client *realsearch.Client) gin.HandlerFunc {
	proxy := client.NewProxy()
	return func(ctx *gin.Context) {
		ctx.Request.URL.Path = ctx.Param("path")
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
