package api

import (
	"errors"
	"fmt"
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/queue"
	"github.com/acgn-org/onest/internal/server/response"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/repository"
	"github.com/acgn-org/onest/tools"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"regexp"
)

func GetItemDownloads(ctx *gin.Context) {
	id, err := tools.IDFromParam(ctx, "id")
	if err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	downloadRepo := database.NewRepository[repository.DownloadRepository]()
	tasks, err := downloadRepo.GetDownloadTaskByID(id)
	if err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	queue.MigrateDownloadTaskInfo(tasks)

	if tasks == nil {
		tasks = make([]repository.DownloadTask, 0)
	}
	response.Success(ctx, tasks)
}

func GetItems(ctx *gin.Context) {
	var form struct {
		ActiveAfter int32 `form:"active_after" json:"active_after" binding:"min=0"`
	}
	if err := ctx.ShouldBind(&form); err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	itemRepo := database.NewRepository[repository.ItemRepository]()
	items, err := itemRepo.GetWithDateEnd(form.ActiveAfter)
	if err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	if items == nil {
		items = make([]repository.Item, 0)
	}
	response.Success(ctx, items)
}

func NewItem(ctx *gin.Context) {
	var form struct {
		repository.NewItemForm
		Downloads []repository.DownloadForm `json:"downloads" form:"downloads"`
	}
	if err := ctx.ShouldBind(&form); err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	_, err := regexp.Compile(form.Regexp)
	if err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	itemRepo := database.BeginRepository[repository.ItemRepository]()
	defer itemRepo.Rollback()

	item, err := itemRepo.CreateWithForm(&form.NewItemForm)
	if err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	var downloadModels = make([]repository.Download, 0, len(form.Downloads))
	for _, download := range form.Downloads {
		msg, err := source.Telegram.GetMessage(item.ChannelID, download.MsgID)
		if err != nil {
			response.Error(ctx, response.ErrTelegram, err)
			return
		}
		messageVideo, ok := source.Telegram.GetMessageVideo(msg)
		if !ok {
			response.ErrorWithTip(ctx, response.ErrTelegram, fmt.Sprintf("message %d is not video message", download.MsgID))
			return
		}
		downloadModels = append(downloadModels, repository.Download{
			ItemID:   item.ID,
			MsgID:    download.MsgID,
			Text:     messageVideo.Caption.Text,
			Size:     messageVideo.Video.Video.Size,
			Date:     msg.Date,
			Priority: download.Priority,
		})
	}
	if len(downloadModels) != 0 {
		downloadRepo := repository.DownloadRepository{Repository: itemRepo.Repository}
		if err := downloadRepo.CreateAll(downloadModels); err != nil {
			response.Error(ctx, response.ErrDBOperation, err)
			return
		}
	}

	if err := itemRepo.Commit().Error; err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	select {
	case queue.ActivateTaskControl <- struct{}{}:
	default:
	}

	response.Default(ctx)
}

func DeleteItem(ctx *gin.Context) {
	var form struct {
		DeleteTargetPath bool `json:"delete_target_path" form:"delete_target_path"`
	}
	if err := ctx.ShouldBind(&form); err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	id, err := tools.IDFromParam(ctx, "id")
	if err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	itemRepo := database.BeginRepository[repository.ItemRepository]()
	defer itemRepo.Rollback()

	_, err = itemRepo.FirstItemByIDForUpdates(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(ctx, response.ErrNotFound)
			return
		}
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	downloadRepo := repository.DownloadRepository{Repository: itemRepo.Repository}

	downloadIDs, err := downloadRepo.GetIDByItemForUpdates(id)
	if err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	queue.RemoveTasks(downloadIDs...)

	if err := itemRepo.Commit().Error; err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	select {
	case queue.ActivateTaskControl <- struct{}{}:
	default:
	}

	response.Default(ctx)
}

func PatchItem(ctx *gin.Context) {
	var form struct {
		Name       string `json:"name" form:"name"`
		Regexp     string `json:"regexp" form:"regexp"`
		Pattern    string `json:"pattern" form:"pattern"`
		TargetPath string `json:"target_path" form:"target_path"`
	}
	if err := ctx.ShouldBind(&form); err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	id, err := tools.IDFromParam(ctx, "id")
	if err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	itemRepo := database.BeginRepository[repository.ItemRepository]()
	defer itemRepo.Rollback()

	ok, err := itemRepo.UpdatesItemByID(&repository.Item{
		ID:         id,
		Name:       form.Name,
		Regexp:     form.Regexp,
		Pattern:    form.Pattern,
		TargetPath: form.TargetPath,
	})
	if err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	} else if !ok {
		response.Error(ctx, response.ErrNotFound)
		return
	}

	if err := itemRepo.Commit().Error; err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	response.Default(ctx)
}
