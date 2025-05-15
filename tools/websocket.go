package tools

import (
	"github.com/gorilla/websocket"
	"time"
)

var WebsocketUpgrader = websocket.Upgrader{
	HandshakeTimeout: time.Minute,
}
