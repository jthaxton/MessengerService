package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"

	// "log"
	"errors"
)

type Handler struct {
	SocketStore *SocketStore
}

type MessageType string

const (
	InstantMessage MessageType = "INSTANT_MESSAGE"
	Notification MessageType = "NOTIFICATION"
	InstantMessageError MessageType = "INSTANT_MESSAGE_ERROR"
)

const AUTH_ENDPOINT = "http://localhost:3000"

type Message struct {
	Type    MessageType `json:"type"`
	Content string `json:"content"`
	SentBy  string `json:"sent_by"`
	SentTo  string `json:"sent_to"`
}

type AuthResponse struct {
	Email   string      `json:"email"`
}

func (handler *Handler) HandlePing(ctx *gin.Context) {
	ctx.JSON(200, map[string]string{"test": "pong"})
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
	// marsh, err := marshal(res.Body)
	if err != nil {
		return nil, err
	}
	// _, err = unmarshal(marsh, &authRes)
	// if err != nil {
	// 	return nil, err
	// }

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
	fmt.Println("RES")
	fmt.Println(res)
	if err != nil {
		fmt.Println(err)
		// ctx.JSON(402, map[string]string{"auth": err.Error()})
		return
	}
	ws, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		fmt.Println(err)
		// ctx.JSON(500, map[string]string{"ws": err.Error()})
	}
	handler.SocketStore.UnsafeAddToStore(res.Email, ws)
	defer ws.Close()

	var messageObj Message
	for {
		//Read Message from client

		mt, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}

		_, err = unmarshal(message, &messageObj)
		fmt.Println(" UNMARSHAlED")

		if err != nil {
			fmt.Println(err)
			// m, err := marshal(Message{Type: InstantMessageError})
			// if err != nil {
			// 	fmt.Println(" UNMARSHAl err")

			// } else {
			// 	fmt.Println("WRITING MESSAGE AFTER UNMARSHAl")
			// 	ws.WriteMessage(mt, m)
			// }
		}
		//If client message is ping will return pong
		// if string(message) == "ping" {
		// 	message = []byte("pong")
		// }
		//Response message to client
		toEmail := messageObj.SentTo
		socket := handler.SocketStore.Sockets[toEmail]
		fmt.Println("socket " + toEmail)
		fmt.Println(socket == nil)
		for k := range handler.SocketStore.Sockets {
			fmt.Println(k)
			// i++
	}
		if socket != nil {
			fmt.Println("WRITING MESSAGE")
			err = socket.WriteMessage(mt, message)
			if err != nil {
				// ctx.JSON(500, map[string]string{"ws": err.Error()})
				break
			}
		}
	}
}