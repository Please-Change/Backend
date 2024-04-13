package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

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

	defer conn.Close()

	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read failed:", err)
			break
		}
		input := string(message)
		message = []byte(input)
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.Println("write failed:", err)
			break
		}
	}
}
