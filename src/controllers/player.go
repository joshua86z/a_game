package controllers

import (
	"code.google.com/p/go.net/websocket"
)

// 玩家
type Player struct {
	Uid        int64
	Token      string
	Conn       *websocket.Conn
	ActiveTime int64
}
