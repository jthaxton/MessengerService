package main

import (
	// "log"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
		r := gin.Default()
		m := make(map[string]*websocket.Conn)
		socketStore := SocketStore{Sockets: m}
		handler := Handler{SocketStore: &socketStore}

		// r.POST("/create", )
		r.GET("/ping", handler.HandlePing)
		r.GET("/connect", handler.HandleConnect)
		r.Run()
}