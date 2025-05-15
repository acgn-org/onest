package response

import (
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func Error(ctx *gin.Context, msg *Msg, args ...any) {
	log.WithFields(log.Fields{
		"code": msg.Code,
		"msg":  msg.Msg,
	}).Debugln(args)
	ctx.AsciiJSON(400, msg)
	ctx.Abort()
}

func ErrorWithTip(ctx *gin.Context, msg *Msg, tip string, args ...any) {
	tipMsg := *msg
	tipMsg.Msg = tip
	Error(ctx, &tipMsg, args...)
}

func Success(ctx *gin.Context, data interface{}) {
	ctx.JSON(200, Msg{
		Data: data,
	})
}

func Default(ctx *gin.Context) {
	Success(ctx, nil)
}
