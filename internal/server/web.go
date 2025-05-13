package server

import (
	"crypto/md5"
	"errors"
	"fmt"
	"github.com/acgn-org/onest/tools"
	"github.com/gin-gonic/gin"
	"io"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
)

func isNotWebRequest(ctx *gin.Context) bool {
	return ctx.Request.Method != "GET" || strings.HasPrefix(ctx.Request.URL.Path, "/api/") || ctx.Writer.Written()
}

func NewWebHandlerWithFS(fe fs.FS) (gin.HandlerFunc, error) {
	file, err := fe.Open("index.html")
	if err != nil {
		return nil, err
	}
	fileContentBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	_ = file.Close()
	index := string(fileContentBytes)
	hash := md5.New()
	hash.Write(fileContentBytes)
	indexEtag := fmt.Sprintf("W/\"%x\"", hash.Sum(nil))

	fileServer := http.StripPrefix("/", http.FileServer(http.FS(fe)))

	return func(ctx *gin.Context) {
		if isNotWebRequest(ctx) {
			return
		}

		f, err := fe.Open(strings.TrimPrefix(ctx.Request.URL.Path, "/"))
		if err != nil {
			var fsError *fs.PathError
			if errors.As(err, &fsError) {
				if ctx.GetHeader("If-None-Match") == indexEtag {
					ctx.AbortWithStatus(304)
					return
				}
				ctx.Header("Content-Type", "text/html")
				ctx.Header("Cache-Control", "no-cache")
				ctx.Header("Etag", indexEtag)
				ctx.String(200, index)
				ctx.Abort()
				return
			}
		}
		_ = f.Close()

		ctx.Header("Cache-Control", "public, max-age=2592000, immutable")
		fileServer.ServeHTTP(ctx.Writer, ctx.Request)
		ctx.Abort()
	}, nil
}

func NewWebHandlerWithAddress(addr string) (gin.HandlerFunc, error) {
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.BufferPool = tools.BufferHttpUtil{}
	proxy.Transport = &http.Transport{}

	return func(ctx *gin.Context) {
		if isNotWebRequest(ctx) {
			return
		}
		proxy.ServeHTTP(ctx.Writer, ctx.Request)
	}, nil
}
