package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/internal/server"
	"github.com/acgn-org/onest/realsearch"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	logger := logfield.New("server")

	_realSearch, err := realsearch.NewClient(&realsearch.Config{
		HttpClient: config.RealSearchHttpClient,
		BaseUrl:    config.RealSearch.BaseUrl,
	})
	if err != nil {
		logger.Fatalln("create real search client failed:", err)
	}

	addr := fmt.Sprintf("%s:%d", config.Server.Host, config.Server.Port)
	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logger.Fatalf("listen on %s failed: %v", addr, err)
	}
	logger.Infof("listening on %s", listen.Addr())

	httpSrv := http.Server{
		Handler: server.New(&server.Config{
			RealSearch: _realSearch,
		}),
	}
	go func(listen net.Listener) {
		if err := httpSrv.Serve(listen); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return
			}
			logger.Fatalln("start server failed:", err)
		}
	}(listen)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt, os.Kill, syscall.SIGTERM)
	<-quit
	logger.Infoln("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Errorln("shutdown http server failed:", err)
	}
}
