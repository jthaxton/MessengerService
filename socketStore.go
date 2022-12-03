package main

import (

	"github.com/gorilla/websocket"
)

type Socket struct {
	Socket     *websocket.Conn
	Kind       int
}