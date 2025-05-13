package main

import (
	"github.com/acgn-org/onest/internal/server"
	log "github.com/sirupsen/logrus"
	"os"
)

func main() {
	engine := server.NewEngine()

	addr := os.Getenv("ONEST_WEB_SERVER_ADDR")
	if addr == "" {
		addr = "http://localhost:5173"
	}
	webHandler, err := server.NewWebHandlerWithAddress(addr)
	if err != nil {
		log.Fatalln("create web handler failed:", err)
	}
	engine.Use(webHandler)

	server.Run(engine)
}
