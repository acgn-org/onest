package api

import (
	"github.com/acgn-org/onest/internal/logfield"
	"github.com/acgn-org/onest/internal/logtee"
	"github.com/acgn-org/onest/internal/server/response"
	"github.com/acgn-org/onest/tools"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

func WatchLogs(ctx *gin.Context) {
	logger := logfield.New(logfield.ComServer).WithAction("websocket").WithField("id", uuid.New().String())

	conn, err := tools.WebsocketUpgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		logger.Warnln("upgrade connection failed:", err)
		response.Error(ctx, response.ErrForm, err)
		return
	}
	defer conn.Close()

	sub := logtee.NewSubscribe()
	defer sub.Close()

	for logs := range sub.Listen() {
		for _, log := range logs {
			if err := conn.WriteMessage(websocket.BinaryMessage, log); err != nil {
				logger.Debugln("failed to write log, exiting:", err)
				return
			}
		}
	}
}
