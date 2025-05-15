package api

import (
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/internal/source"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RealSearchProxy() gin.HandlerFunc {
	logger := logfield.New(logfield.ComServer).WithAction("real search proxy")
	proxy := source.RealSearch.NewProxy()
	proxy.ErrorHandler = func(_ http.ResponseWriter, request *http.Request, err error) {
		logger.WithField("url", request.URL.String()).Warnln("proxy error:", err)
	}
	return func(ctx *gin.Context) {
		ctx.Request.URL.Path = ctx.Param("path")
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}
}
