package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Please-Change/backend/pkg/types"
)

const PORT = 5174

var QueueGroup = new(sync.WaitGroup)
var ReadyPlayerCount = make(chan int64)
var MyGameState = types.GameState{
	Status: types.Pending,
	Settings: types.GameSettings{
		Language: "javascript",
		Problem:  "FizzBuzz",
	},
}

var nextId = make(chan int64)
var _PlayerStateBuffer = map[int64]*types.PlayerState{}
var lock = sync.RWMutex{}

func AddPlayerState(ps *types.PlayerState) int64 {
	lock.Lock()
	defer lock.Unlock()
	// get an id
	var currentId = <-nextId
	_PlayerStateBuffer[currentId] = ps
	nextId <- currentId + 1
	return currentId
}

func BroadcastToAllPlayers(action types.Action, data interface{}) {
	lock.RLock()
	defer lock.RUnlock()
	for id := range _PlayerStateBuffer {
		ps := GetPlayerState(id)
		if ps.Status == types.Ready {
			ps.SendMessage(action, data)
		}
	}
}

func RemovePlayerState(id int64) {
	lock.Lock()
	defer lock.Unlock()
	delete(_PlayerStateBuffer, id)
}

func GetPlayerState(id int64) *types.PlayerState {
	lock.RLock()
	defer lock.RUnlock()
	return _PlayerStateBuffer[id]
}

func Serve() {
	QueueGroup.Add(MAX_QUEUED_PLAYERS)
	fmt.Printf("Serving on %d\n", PORT)

	http.HandleFunc("/game", HandleStart)

	http.HandleFunc("/complete", HandleComplete)

	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
