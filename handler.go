package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	"errors"
)

type Handler struct {
	SocketStore map[string]*Socket
}

type MessageType string

const (
	Conversation MessageType =      "Conversation"
	Notification MessageType =        "NOTIFICATION"
	InstantMessageError MessageType = "INSTANT_MESSAGE_ERROR"
)

const AUTH_ENDPOINT = "http://localhost:3000"

type Message struct {
	Type    MessageType				`json:"data_object_type"`
	Content string      			`json:"content"`
	SentBy  string						`json:"sent_by"`
	SentTo  string      			`json:"sent_to"`
	DataObjectId  string      `json:"data_object_id"`
}

type AuthResponse struct {
	Email   string      `json:"email"`
}

func Authenticate(ctx *gin.Context) (*AuthResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", AUTH_ENDPOINT + "/sessions/get_token", nil)
	if err != nil {
		return nil, err
	}

	token := ctx.Request.URL.Query()["token"]
	req.Header.Set("Authorization", token[0])
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	authRes := AuthResponse{}
	err = json.NewDecoder(res.Body).Decode(&authRes)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		return nil, errors.New("Unauthenticated")
	}

	return &authRes, nil
}

func (handler *Handler) HandleConnect(ctx *gin.Context) {

	var upgrader = websocket.Upgrader{
    //check origin will check the cross region source (note : please not using in production)
		CheckOrigin: func(r *http.Request) bool {
					//Here we just allow the chrome extension client accessable (you should check this verify accourding your client source)
			return true //origin == "chrome-extension://cbcbkhdmedgianpaifchdaddpnmgnknn"
		},
	}

	res, err := Authenticate(ctx)

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err.Error())
	}
	defer ws.Close()
	
	var messageObj Message
	for {
		mt, message, err := ws.ReadMessage()

		if err != nil {
			fmt.Println(err.Error())
			break
		}
		handler.SocketStore[res.Email] = &Socket{Socket: ws, Kind: mt}
		_, err = unmarshal(message, &messageObj)

		if err != nil {
			fmt.Println(err.Error())
		}

		//Response message to client
		toEmail := messageObj.SentTo
		sockets := handler.SocketStore
		socket := sockets[toEmail]
		if socket != nil {
			err = socket.Socket.WriteMessage(mt, message)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
	}
}

func (handler *Handler) HandleDisconnect(ctx *gin.Context) {
	email := ctx.Request.URL.Query()["email"][0]
	delete(handler.SocketStore, email)
	ctx.JSON(200, make(map[string]string))
}

func (handler *Handler) HandleSendMessage(ctx *gin.Context) {
	var messageObj Message
	err := json.NewDecoder(ctx.Request.Body).Decode(&messageObj)

	if err != nil {
		fmt.Println(err.Error())
	}
	toEmail := messageObj.SentTo
	sockets := handler.SocketStore
	socket := sockets[toEmail]
	fmt.Println("toEmail")
	fmt.Println(toEmail)

	// *********************************
	// * TODO validate messageObj here.*
	// *********************************


	message, err := marshal(messageObj)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println("GOT HERE 11111")

	keys := make([]string, len(handler.SocketStore))

	i := 0
	for k := range handler.SocketStore {
			keys[i] = k
			i++
	}
	fmt.Println(keys)

	if socket != nil {
		fmt.Println(message)
		err = socket.Socket.WriteMessage(socket.Kind, message)
		if err != nil {
			fmt.Println(err.Error())
		}
	}
}