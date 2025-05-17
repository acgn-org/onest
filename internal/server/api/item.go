package api

import (
	"github.com/acgn-org/onest/internal/database"
	"github.com/acgn-org/onest/internal/server/response"
	"github.com/acgn-org/onest/repository"
	"github.com/gin-gonic/gin"
)

func GetItems(ctx *gin.Context) {
	var form struct {
		ActiveAfter uint16 `form:"active_after" json:"active_after"`
	}
	if err := ctx.ShouldBind(&form); err != nil {
		response.Error(ctx, response.ErrForm, err)
		return
	}

	itemRepo := database.NewRepository[repository.ItemRepository]()
	items, err := itemRepo.GetWithDateEnd(int32(form.ActiveAfter))
	if err != nil {
		response.Error(ctx, response.ErrDBOperation, err)
		return
	}

	if items == nil {
		items = make([]repository.Item, 0)
	}
	response.Success(ctx, items)
}
