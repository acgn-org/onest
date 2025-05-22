package server

import (
	"github.com/acgn-org/onest/internal/server/api"
	"github.com/gin-gonic/gin"
)

func Api(group *gin.RouterGroup) {
	group.Any("realsearch/*path", api.RealSearchProxy())

	item := group.Group("item")
	item.GET("active", api.GetActiveItems)
	item.GET("error", api.GetErrorItems)
	item.POST("/", api.NewItem)
	itemWithId := item.Group(":id")
	itemWithId.GET("downloads", api.GetItemDownloads)
	itemWithId.PATCH("/", api.PatchItem)
	itemWithId.DELETE("/", api.DeleteItem)

	download := group.Group("download")
	download.GET("tasks", api.GetDownloadTasks)
	downloadWithId := download.Group(":id")
	downloadWithId.PATCH("priority", api.UpdateDownloadPriority)
	downloadWithId.DELETE("/", api.DeleteDownload)

	log := group.Group("log")
	log.GET("watch", api.WatchLogs)
}
