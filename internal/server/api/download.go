package api

import (
	"errors"
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/queue"
	"github.com/acgn-org/onest/internal/server/response"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
	"github.com/acgn-org/onest/tools"
	"github.com/gin-gonic/gin"
	"github.com/zelenin/go-tdlib/client"
	"gorm.io/gorm"
)

func AddDownloadForItem(ctx *gin.Context) {
	var form struct {
		ItemID    uint  `json:"item_id" form:"item_id" binding:"required"`
		MessageID int64 `json:"message_id" form:"message_id" binding:"required"`
		Priority  int32 `json:"priority" form:"priority" binding:"min=1,max=32"`
	}
	if err := ctx.ShouldBind(&form); err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	itemRepo := database.BeginRepository[repository.ItemRepository]()
	defer itemRepo.Rollback()

	item, err := itemRepo.FirstItemByIDForUpdates(form.ItemID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.ErrorWithTip(ctx, response.ErrNotFound, "item does not exist")
			return
		}
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	msg, err := source.Telegram.GetMessage(item.ChannelID, form.MessageID)
	if err != nil {
		response.Error(ctx, response.ErrTelegram, err)
		return
	}

	downloadRepo := repository.DownloadRepository{Repository: itemRepo.Repository}
	result, err := downloadRepo.CreateWithMessages(item.ID, item.Priority, []*client.Message{msg})
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			response.Error(ctx, response.ErrResourceConflict, "message already exists")
			return
		}
		response.Error(ctx, response.ErrDBOperation, err)
		return
	} else if len(result) == 0 {
		response.ErrorWithTip(ctx, response.ErrTelegram, "message dose not contain video")
		return
	}

	if err := itemRepo.Commit().Error; err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	queue.TryActivateTaskControl()

	response.Default(ctx)
}

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

func ForceStartTask(ctx *gin.Context) {
	id, err := tools.UintIDFromParam(ctx, "id")
	if err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	downloadRepo := database.NewRepository[repository.DownloadRepository]()
	downloadTask, err := downloadRepo.FirstByIDPreloadItem(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(ctx, response.ErrNotFound)
			return
		}
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	if err := queue.ForceAddDownloadQueue(downloadTask.Item.ChannelID, *downloadTask); err != nil {
		response.Error(ctx, response.ErrUnexpected, err)
		return
	}

	response.Default(ctx)
}

func ForceResetTask(ctx *gin.Context) {
	id, err := tools.UintIDFromParam(ctx, "id")
	if err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	downloadRepo := database.BeginRepository[repository.DownloadRepository]()
	defer downloadRepo.Rollback()

	ok, err := downloadRepo.UpdateResetDownloadState(id)
	if err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	} else if !ok {
		response.Error(ctx, response.ErrNotFound)
		return
	}

	queue.RemoveTasks(id)

	if err := downloadRepo.Commit().Error; err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	queue.TryActivateTaskControl()

	response.Default(ctx)
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
