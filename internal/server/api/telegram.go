package api

import (
	"context"
	"time"

	"github.com/acgn-org/onest/internal/config"
	"github.com/acgn-org/onest/internal/server/response"
	"github.com/acgn-org/onest/internal/source"
	"github.com/acgn-org/onest/tools"
	"github.com/gin-gonic/gin"
)

func GetChat(ctx *gin.Context) {
	id, err := tools.Int64IDFromParam(ctx, "id")
	if err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	_ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(config.Server.Get().Timeout))
	defer cancel()

	info, err := source.Telegram.GetChat(_ctx, id)
	if err != nil {
		response.Error(ctx, response.ErrTelegram, err)
		return
	}

	response.Success(ctx, info)
}

func GetChatPhoto(ctx *gin.Context) {
	id, err := tools.Int64IDFromParam(ctx, "id")
	if err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	_ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(config.Server.Get().Timeout))
	defer cancel()

	info, err := source.Telegram.GetChat(_ctx, id)
	if err != nil {
		response.Error(ctx, response.ErrTelegram, err)
		return
	}

	if !info.Photo.Big.Local.IsDownloadingCompleted {
		file, err := source.Telegram.DownloadFile(info.Photo.Big.Id, 32, true)
		if err != nil {
			response.Error(ctx, response.ErrTelegram, err)
			return
		}
		info.Photo.Big.Local = file.Local
	}

	ctx.File(info.Photo.Big.Local.Path)
}

func GetMessage(ctx *gin.Context) {
	chatID, err := tools.Int64IDFromParam(ctx, "id")
	if err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	messageID, err := tools.Int64IDFromParam(ctx, "msgId")
	if err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	_ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(config.Server.Get().Timeout))
	defer cancel()

	info, err := source.Telegram.GetMessage(_ctx, chatID, messageID)
	if err != nil {
		response.Error(ctx, response.ErrTelegram, err)
		return
	}

	response.Success(ctx, info)
}
