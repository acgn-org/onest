package tools

import (
	"time"

	"github.com/gorilla/websocket"
)

var WebsocketUpgrader = websocket.Upgrader{
	HandshakeTimeout: time.Minute,
}
