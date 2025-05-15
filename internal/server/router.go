package server

import (
	"github.com/acgn-org/onest/internal/server/api"
	"github.com/gin-gonic/gin"
)

func Api(group *gin.RouterGroup) {
	group.Any("realsearch/*path", api.RealSearchProxy())

	download := group.Group("download")
	download.GET("tasks", api.GetDownloadTasks)

	log := group.Group("log")
	log.GET("watch", api.WatchLogs)
}
