package server

import (
	"bytes"
	"log"
	"net/http"
	"sync"

	"github.com/Please-Change/backend/pkg/types"
	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
)

const MAX_QUEUED_PLAYERS = 128

func ProcessGame(gs *types.GameState) {
	QueueGroup.Add(1)
	defer gs.Socket.Close()
	defer QueueGroup.Done()
	for {
		mt, message, err := gs.Socket.ReadMessage()
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

					ReadyPlayerCount <- <-ReadyPlayerCount + 1
					var msg = map[string]interface{}{
						"action": types.Action(types.PlayerCount),
						"data":   <-ReadyPlayerCount,
					}

					var w = bytes.NewBuffer(nil)
					var enc = sonic.ConfigDefault.NewEncoder(w)
					enc.Encode(msg)
					gs.Socket.WriteMessage(mt, w.Bytes())
				}
			}
		case types.ChangeSetting:
			{
				if gs.Ready == types.ReadyState(types.Waiting) {

				}
			}
		case types.Submit:
			{
				if gs.Ready == types.Active {
					// Check
				}
			}
		case types.StatusRequest:
			{

			}
		}

	}
}

func HandleStart(w http.ResponseWriter, r *http.Request) {

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
	QueueGroup.Wait()
}

func HandleComplete(w http.ResponseWriter, r *http.Request) {

}
