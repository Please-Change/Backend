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

func ProcessGame(ps *types.PlayerState) {
	QueueGroup.Add(1)
	defer ps.Socket.Close()
	defer QueueGroup.Done()
	for {
		mt, message, err := ps.Socket.ReadMessage()
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
			break
		}

		log.Printf("Action: %s\n", action)

		switch action {
		case types.UsePowerUp:
			{
				if ps.Status == types.Active {
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
				if ps.Status == types.ReadyState(types.Ready) {
					ps.Status = types.ReadyState(types.Waiting)
				} else if ps.Status == types.ReadyState(types.Waiting) {
					ps.Status = types.ReadyState(types.Ready)

					ReadyPlayerCount <- <-ReadyPlayerCount + 1
					var msg = map[string]interface{}{
						"action": types.Action(types.PlayerCount),
						"data":   <-ReadyPlayerCount,
					}

					var w = bytes.NewBuffer(nil)
					var enc = sonic.ConfigDefault.NewEncoder(w)
					enc.Encode(msg)
					ps.Socket.WriteMessage(mt, w.Bytes())
				}
			}
		case types.ChangeSetting:
			{
				if ps.Status == types.ReadyState(types.Waiting) {

				}
			}
		case types.Submit:
			{
				if ps.Status == types.Active {
					// Check
				}
			}
		case types.StatusRequest:
			{
				// var msg = map[string]interface{}{
				// 	"action": types.Action(types.PlayerCount),
				// 	"data": map[string]interface{}{
				// 		"status": ps.Status,
				// 	},
				// }

				// var w = bytes.NewBuffer(nil)
				// var enc = sonic.ConfigDefault.NewEncoder(w)
				// enc.Encode(msg)
				// gs.Socket.WriteMessage(mt, w.Bytes())

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

	game := types.PlayerState{
		Status: types.Waiting,
		// Settings: types.GameSettings{
		// 	Language: "javascript",
		// 	Problem:  "FizzBuzz",
		// },
		Socket: conn,
	}

	go ProcessGame(&game)
	QueueGroup.Wait()
}

func HandleComplete(w http.ResponseWriter, r *http.Request) {

}
