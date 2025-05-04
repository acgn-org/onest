package server

import (
	"github.com/acgn-org/onest/realsearch"
	"github.com/gin-gonic/gin"
)

type Config struct {
	RealSearch *realsearch.Client
}

func New(c *Config) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	Engine := gin.Default()

	Api(c, Engine.Group("api"))

	return Engine
}
