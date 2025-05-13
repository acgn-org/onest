package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/gin-gonic/gin"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func NewEngine() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	Engine := gin.Default()

	Api(Engine.Group("api"))

	return Engine
}

func Run(engine *gin.Engine) {
	logger := logfield.New(logfield.ComServer)
	serverConfig := config.Server.Get()

	// create http server

	addr := fmt.Sprintf("%s:%d", serverConfig.Host, serverConfig.Port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalf("listen on %s failed: %v", addr, err)
	}
	logger.Infof("listening on %s", listen.Addr())

	httpSrv := http.Server{
		Handler: engine,
	}
	go func(listen net.Listener) {
		if err := httpSrv.Serve(listen); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			logger.Fatalln("start server failed:", err)
		}
	}(listen)

	// shutdown

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-quit
	logger.Infoln("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		logger.Errorln("shutdown http server failed:", err)
	}
}
