package main

import (
	"errors"
	"fmt"

	"github.com/gorilla/websocket"
)


type Sockets map[string]*websocket.Conn
type SocketStore struct {
	Sockets Sockets
}

func (store *SocketStore) UnsafeAddToStore(email string, value *websocket.Conn) {
	fmt.Println("ADDING TO STORE")
	// fmt.Println(value)
	fmt.Println(email)
	store.Sockets[email] = value
}

func (store *SocketStore) SafeAddToStore(email string, value *websocket.Conn) error {
	if store.Sockets[email] != nil {
		store.Sockets[email] = value
		return nil
	}

	return errors.New("socket for email already exists")
}