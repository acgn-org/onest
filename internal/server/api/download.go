package api

import (
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/queue"
	"github.com/acgn-org/onest/internal/server/response"
	"github.com/acgn-org/onest/repository"
	"github.com/acgn-org/onest/tools"
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

func DeleteDownload(ctx *gin.Context) {
	id, err := tools.IDFromParam(ctx, "id")
	if err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	downloadRepo := database.BeginRepository[repository.DownloadRepository]()
	defer downloadRepo.Rollback()

	if err := downloadRepo.DeleteByID(id); err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	if err := downloadRepo.Commit().Error; err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	queue.RemoveTasks(id)

	response.Default(ctx)
}
