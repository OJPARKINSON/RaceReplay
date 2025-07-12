package main

import (
	"net/http"

	"golang.org/x/net/websocket"
)

func MessageServer(ws *websocket.Conn) {
	println(ws.Request().URL.Query().Get("sessionID"))
	defer ws.Close()
	ws.Write([]byte("hello"))
	ws.Write([]byte("ollie"))
}

func main() {
	http.Handle("/hello", websocket.Handler(MessageServer))
	err := http.ListenAndServe(":8765", nil)
	if err != nil {
		panic("ListenAndServe: " + err.Error())
	}
}
