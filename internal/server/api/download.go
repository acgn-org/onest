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
	if tasks == nil {
		tasks = make([]repository.DownloadTask, 0)
	}
	response.Success(ctx, tasks)
}

func UpdateDownloadPriority(ctx *gin.Context) {
	var form struct {
		Priority int32 `json:"priority" form:"priority" binding:"min=1,max=32"`
	}
	if err := ctx.ShouldBind(&form); err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	id, err := tools.UintIDFromParam(ctx, "id")
	if err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	downloadRepo := database.BeginRepository[repository.DownloadRepository]()
	defer downloadRepo.Rollback()

	ok, err := downloadRepo.UpdatePriority(id, form.Priority)
	if err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	} else if !ok {
		response.Error(ctx, response.ErrNotFound)
		return
	}

	if err := downloadRepo.Commit().Error; err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	queue.UpdatePriority(id, form.Priority)

	response.Default(ctx)
}

func DeleteDownload(ctx *gin.Context) {
	id, err := tools.UintIDFromParam(ctx, "id")
	if err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	downloadRepo := database.BeginRepository[repository.DownloadRepository]()
	defer downloadRepo.Rollback()

	ok, err := downloadRepo.DeleteByID(id)
	if err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	} else if !ok {
		response.Error(ctx, response.ErrNotFound)
		return
	}

	if err := downloadRepo.Commit().Error; err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	queue.RemoveTasks(id)

	response.Default(ctx)
}
