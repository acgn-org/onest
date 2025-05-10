package source

import (
	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/realsearch"
)

var RealSearch *realsearch.Client

func init() {
	var err error
	RealSearch, err = realsearch.NewClient(&realsearch.Config{
		HttpClient: config.RealSearchHttpClient,
		BaseUrl:    config.RealSearch.Get().BaseUrl,
	})
	if err != nil {
		logfield.New(logfield.ComSource).WithAction("init:realsearch").
			Fatalln("create real search client failed:", err)
	}
}
