package main

import (
	"net/http"
	"ojparkinson/RaceReplay/internal/session"

	"golang.org/x/net/websocket"
)

func MessageServer(ws *websocket.Conn) {
	sessionID := ws.Request().URL.Query().Get("sessionID")
	sessionReplay := session.NewSessionReplay("session123")
	sessionReplay.SpeedMultiplier = 1.5

	sessionReplay.Start(ws)

	defer ws.Close()
	ws.Write([]byte("Starting stream for session " + sessionID))

}

func main() {
	http.Handle("/hello", websocket.Handler(MessageServer))
	err := http.ListenAndServe(":8765", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
