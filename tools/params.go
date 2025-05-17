package tools

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func IDFromParam(ctx *gin.Context, name string) (uint, error) {
	id64, err := strconv.ParseUint(ctx.Param(name), 10, 64)
	if err != nil {
		return 0, err
	}
	if id64 == 0 {
		return 0, fmt.Errorf("id is empty")
	}
	return uint(id64), nil
}
