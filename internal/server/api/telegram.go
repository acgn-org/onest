package api

import (
	"github.com/acgn-org/onest/internal/server/response"
	"github.com/acgn-org/onest/internal/source"
	"github.com/gin-gonic/gin"
	"strconv"
)

func GetChat(ctx *gin.Context) {
	idStr := ctx.Param("id")
	if idStr == "" {
		response.Error(ctx, response.ErrForm)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	info, err := source.Telegram.GetChat(id)
	if err != nil {
		response.Error(ctx, response.ErrTelegram, err)
		return
	}

	response.Success(ctx, info)
}
