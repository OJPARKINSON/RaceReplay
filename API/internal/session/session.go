package session

import (
	"golang.org/x/net/websocket"
)

type SessionReplay struct {
	SessionID string
}

func (s *SessionReplay) Start(ws *websocket.Conn) {
	println(s.SessionID)
}
