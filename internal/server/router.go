package server

import (
	"github.com/acgn-org/onest/internal/server/api"
	"github.com/gin-gonic/gin"
)

func Api(group *gin.RouterGroup) {
	group.Any("realsearch/*path", api.RealSearchProxy())

	item := group.Group("item")
	item.GET("active", api.GetItems)
	item.POST("/", api.NewItem)
	item.DELETE("/:id", api.DeleteItem)

	download := group.Group("download")
	download.GET("tasks", api.GetDownloadTasks)

	log := group.Group("log")
	log.GET("watch", api.WatchLogs)
}
