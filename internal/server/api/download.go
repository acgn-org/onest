package api

import (
	"github.com/acgn-org/onest/internal/queue"
	"github.com/acgn-org/onest/internal/server/response"
	"github.com/gin-gonic/gin"
)

func GetDownloadTasks(ctx *gin.Context) {
	tasks, err := queue.GetDownloading()
	if err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}
	response.Success(ctx, tasks)
}
