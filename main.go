package main

import (

	"github.com/gin-gonic/gin"
)

func main() {
		r := gin.Default()
		r.UnescapePathValues = false

		socketStore := make(map[string]*Socket)
		
		handler := Handler{SocketStore: socketStore}

		r.GET("/connect", handler.HandleConnect)
		r.POST("/disconnect", handler.HandleDisconnect)
		r.POST("/send_message", handler.HandleSendMessage)
		r.Run()
}