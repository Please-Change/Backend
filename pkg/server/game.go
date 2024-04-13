package server

import (
	"github.com/Please-Change/backend/pkg/types"
	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

func ProcessGame(conn *websocket.Conn, gs *types.GameState) {
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

		switch action {
		case types.KeyPress:
			{
				key, err := root.Get("key").String()
				if err != nil {
					log.Printf("read: %s", err)
				}
				log.Println("Entered: ", key)

			}
		case types.UsePowerUp:
			{
			}
		case types.TogglePause:
			{
			}
		case types.QuitGame:
			{
				break
			}
		}

	}
	defer conn.Close()
}

func HandleStart(w http.ResponseWriter, r *http.Request) {
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

	newGame := types.GameState{}

	go ProcessGame(conn, &newGame)
}
