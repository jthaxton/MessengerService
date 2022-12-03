package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func main() {
		r := gin.Default()
		r.UnescapePathValues = false
		AUTH_ENDPOINT := os.Getenv("RAILS_APP")
		socketStore := make(map[string]*Socket)
		
		handler := Handler{SocketStore: &socketStore, AuthEndpoint: AUTH_ENDPOINT}

		r.GET("/connect", handler.HandleConnect)
		r.POST("/disconnect", handler.HandleDisconnect)
		r.POST("/send_message", handler.HandleSendMessage)
		r.Run()
}