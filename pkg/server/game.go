package server

import (
	"log"
	"net/http"
	"sync"

	"github.com/Please-Change/backend/pkg/types"
	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
)

const MAX_SIMULTANEOUS_PLAYERS = 8
const MAX_QUEUED_PLAYERS = 128

func ProcessGame(gs *types.GameState) {
	defer gs.Socket.Close()
	for {
		_, message, err := gs.Socket.ReadMessage()
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

		log.Printf("Action: %s\n", action)

		switch action {
		case types.UsePowerUp:
			{
				if gs.Ready == types.Active {
					// Use the power up
					used, err := root.Get("use").Int64()
					if err != nil {
						log.Printf("read: %s\n", err)
					}
					log.Printf("Used item %d\n", used)
				}
			}
		case types.ChangeReady:
			{
				if gs.Ready == types.ReadyState(types.Ready) {
					gs.Ready = types.ReadyState(types.Waiting)
				} else if gs.Ready == types.ReadyState(types.Waiting) {
					gs.Ready = types.ReadyState(types.Ready)
				}
			}
		case types.ChangeSetting:
			{
				if gs.Ready == types.Waiting {

				}
			}
		case types.Submit:
			{
				if gs.Ready == types.Active {

				}
			}
		case types.StatusRequest:
			{

			}
		}

	}
}

func HandleStart(w http.ResponseWriter, r *http.Request) {
	if len(queue) >= MAX_QUEUED_PLAYERS {
		log.Println("Queued Player Limit Exceeded!")
		return
	}
	var upgrade = websocket.Upgrader{
		ReadBufferSize:  512,
		WriteBufferSize: 512,
		WriteBufferPool: &sync.Pool{},
		CheckOrigin:     func(r *http.Request) bool { return true },
	}

	conn, err := upgrade.Upgrade(w, r, nil)

	if err != nil {
		log.Println("Couldn't upgrade connection: ", err)
		return
	}

	game := types.GameState{
		Ready:    types.Waiting,
		Settings: types.GameSettings{},
		Socket:   conn,
	}

	go ProcessGame(&game)
}

func HandleComplete(w http.ResponseWriter, r *http.Request) {

}
