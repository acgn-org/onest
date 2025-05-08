package server

import (
	"github.com/acgn-org/onest/internal/server/api"
	"github.com/gin-gonic/gin"
)

func Api(group *gin.RouterGroup) {
	group.Any("realsearch/*path", api.RealSearchProxy())
}
