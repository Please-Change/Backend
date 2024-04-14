package server

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Please-Change/backend/pkg/types"
)

const PORT = 5174

var MyGameState = types.GameState{
	Status: types.Pending,
	Settings: types.GameSettings{
		Language: "javascript",
		Problem:  "FizzBuzz",
	},
}

type PlayerStateStore struct {
	sync.RWMutex
	store  map[int64]*types.PlayerState
	nextId int64
}

var Players PlayerStateStore

func (sb *PlayerStateStore) Add(ps *types.PlayerState) int64 {
	sb.Lock()
	defer sb.Unlock()
	// get an id
	var currentId = sb.nextId
	sb.store[currentId] = ps
	sb.nextId = currentId + 1
	return currentId
}

func (sb *PlayerStateStore) BroadcastWithSkip(mt int, action types.Action, data interface{}, skip int64) {
	sb.RLock()
	defer sb.RUnlock()
	for id := range sb.store {
		if id == skip {
			continue
		}

		ps := sb.Get(id)
		if ps.Status == types.Ready || ps.Status == types.Active {
			ps.SendMessage(mt, action, data)
		}
	}
}

func (sb *PlayerStateStore) Broadcast(mt int, action types.Action, data interface{}) {
	sb.BroadcastWithSkip(mt, action, data, -1)
}

func (sb *PlayerStateStore) Remove(id int64) {
	sb.Lock()
	defer sb.Unlock()
	delete(sb.store, id)
}

func (sb *PlayerStateStore) Get(id int64) *types.PlayerState {
	sb.RLock()
	defer sb.RUnlock()
	return sb.store[id]
}

func (sb *PlayerStateStore) UpdateStatusFor(id int64, s types.ReadyState) {
	sb.RLock()
	defer sb.RUnlock()
	sb.store[id].SafeSetStatus(s)
}

func (sb *PlayerStateStore) CountReady() int {
	sb.RLock()
	defer sb.RUnlock()

	count := 0
	for id := range sb.store {

		ps := sb.Get(id)
		if ps.Status == types.Ready {
			count++
		}
	}

	return count
}

func Serve() {
	Players = PlayerStateStore{
		store:  map[int64]*types.PlayerState{},
		nextId: 0,
	}
	fmt.Printf("Serving on %d\n", PORT)

	http.HandleFunc("/game", HandleStart)

	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
