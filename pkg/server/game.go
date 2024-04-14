package server

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/Please-Change/backend/pkg/types"
	"github.com/bytedance/sonic"
	"github.com/gorilla/websocket"
)

const MAX_QUEUED_PLAYERS = 128

func ProcessGame(id int64) {
	defer func() {
		ps := Players.Get(id)
		ps.Socket.Close()
	}()

	for {
		ps := Players.Get(id)
		mt, message, err := ps.Socket.ReadMessage()

		log.Printf("[%d] Received %s\n", id, string(message))

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

		switch action {
		case types.UsePowerUp:
			{
				if ps.Status == types.Active {
					// Use the power up
					used, err := root.Get("data").String()
					if err != nil {
						log.Printf("read: %s\n", err)
					}

					Players.BroadcastWithSkip(mt, types.Action(types.UsePowerUp), used, id)
				}
			}
		case types.ChangeReady:
			{
				status, err := root.Get("data").String()
				if err != nil {
					log.Printf("read: %s\n", err)
				}

				if ps.Status != types.ReadyState(status) {
					Players.UpdateStatusFor(id, types.ReadyState(status))

					if status == types.Ready {
						Players.Broadcast(mt, types.Action(types.PlayerCount), Players.CountReady())
						ps.SendMessage(mt, types.ChangeSetting,
							map[string]interface{}{
								"language": MyGameState.Settings.Language,
								"problem":  MyGameState.Settings.Problem,
							},
						)
					}
				}

				ps.SendMessage(mt, types.ChangeReady, status)
			}
		case types.ChangeSetting:
			if ps.Status == types.ReadyState(types.Ready) && MyGameState.Status == types.Pending {
				language, err := root.Get("data").Get("language").String()
				if err != nil {
					log.Printf("read: %s\n", err)
					break
				}

				problem, err := root.Get("data").Get("problem").String()
				if err != nil {
					log.Printf("read: %s\n", err)
					break
				}

				MyGameState.SafeSetSettings(types.GameSettings{
					Language: types.Language(language),
					Problem:  problem,
				})

				Players.Broadcast(mt, types.Action(types.ChangeSetting), map[string]interface{}{
					"language": language,
					"problem":  problem,
				})
			}
		case types.StatusChanged:
			{
				if ps.Status == types.ReadyState(types.Ready) && MyGameState.Status == types.Pending {
					state, err := root.Get("data").Get("status").String()
					if err != nil {
						log.Printf("read: %s\n", err)
						break
					}

					if state != types.Running {
						log.Printf("Incorrect state, cannot start")
						break
					}

					log.Printf("Updating game status %s", types.Running)

					MyGameState.SafeSetStatus(types.Running)

					Players.Broadcast(mt, types.ChangeReady, types.Active)
					Players.Broadcast(mt, types.Action(types.StatusChanged), map[string]interface{}{
						"status": types.Running,
					})
					Players.UpdateStatus(types.Active)
				}
			}
		case types.Submit:
			{
				if ps.Status == types.Active {
					// Check
					program, err := root.Get("data").Get("program").String()
					if err != nil {
						log.Printf("read: %s", err)
					}
					log.Printf("Program, %s\n", program)

					language, err := root.Get("data").Get("language").String()
					if err != nil {
						log.Printf("read: %s", err)
					}

					log.Printf("Language, %s\n", language)

					go processSubmition(id, mt, program, types.Language(language), MyGameState.Settings.Problem)
				} else {
					log.Printf("Incorrect status %s\n", ps.Status)
				}
			}
		case types.StatusRequest:
			{
				var msg = map[string]interface{}{
					"action": types.Action(types.PlayerCount),
					"data": map[string]interface{}{
						"status": MyGameState.Status,
					},
				}

				var w = bytes.NewBuffer(nil)
				var enc = sonic.ConfigDefault.NewEncoder(w)
				enc.Encode(msg)
				ps.Socket.WriteMessage(mt, w.Bytes())
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

	ps := types.PlayerState{
		Status: types.Waiting,
		Socket: conn,
	}

	id := Players.Add(&ps)

	log.Printf("Added player %d\n", id)
	go ProcessGame(id)
}

func processSubmition(id int64, mt int, program string, language types.Language, challenge string) {
	examiner := NewExaminer()

	content, err := os.ReadFile("challenges/" + challenge + ".txt")
	if err != nil {
		log.Fatal(err)
		return
	}

	expectedOutput := string(content)
	res := examiner.RunExam(program, language, expectedOutput)

	ps := Players.Get(id)
	if MyGameState.Status == types.Running {
		if res == "" {
			// Secretly split to waiting
			MyGameState.SetStatus(types.Pending)
			Players.BroadcastWithSkip(mt,
				types.StatusChanged,
				map[string]interface{}{
					"status":  types.End,
					"success": false,
				},
				id,
			)
			ps.SendMessage(mt, types.StatusChanged,
				map[string]interface{}{
					"status":  types.End,
					"success": true,
				},
			)

			// Secretly update player status
			Players.UpdateStatus(types.Waiting)
		} else {
			ps.SendMessage(mt, types.SubmitFailed, res)
		}
	}
}
