package server

import (
	"github.com/acgn-org/onest/internal/server/api"
	"github.com/gin-gonic/gin"
)

func Api(c *Config, group *gin.RouterGroup) {
	group.Any("realsearch/*path", api.RealSearchProxy(c.RealSearch))
}
