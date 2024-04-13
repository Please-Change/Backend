package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
)

func Process(conn *websocket.Conn) {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Printf("read: %s\n", err)
			break
		}
		root, err := sonic.GetFromString(string(message))
		if err != nil {
			log.Printf("read: %s\n", err)
			break
		}
		action, err := root.Get("action").String()
		if err != nil {
			log.Printf("read: %s\n", err)
		}
		log.Printf("recv: %s\n", action)
	}
	defer conn.Close()
}

func handleGame(w http.ResponseWriter, r *http.Request) {
	var upgrade = websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 512,
		WriteBufferPool: &sync.Pool{},
	}

	upgrade.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := upgrade.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Couldn't upgrade connection: ", err)
		return
	}

	go Process(conn)
}
