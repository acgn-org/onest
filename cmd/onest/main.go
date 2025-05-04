package main

import (
	"github.com/acgn-org/onest/internal/server"
	"github.com/acgn-org/onest/realsearch"
	log "github.com/sirupsen/logrus"
)

func main() {
	_realSearch, err := realsearch.NewClient(&realsearch.Config{
		HttpClient: nil, // todo
		BaseUrl:    "",  // todo
	})
	if err != nil {
		log.Fatalln("create real search client failed:", err)
	}

	httpEngine := server.New(&server.Config{
		RealSearch: _realSearch,
	})
	// todo set listen addr
	if err := httpEngine.Run(":80"); err != nil {
		log.Fatalln("run http server on 0.0.0.0:80 failed:", err)
	}
}
