package main

import (
	"github.com/acgn-org/onest/internal/server"
	"github.com/acgn-org/onest/web"
	log "github.com/sirupsen/logrus"
)

func main() {
	engine := server.NewEngine()

	fs, err := web.Fs()
	if err != nil {
		log.Fatal("load web embed files failed, please ensure web/dist existing:", err)
	}
	webHandler, err := server.NewWebHandlerWithFS(fs)
	if err != nil {
		log.Fatalln("create web handler failed:", err)
	}
	engine.Use(webHandler)

	server.Run(engine)
}
