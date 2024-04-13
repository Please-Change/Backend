package types

import (
	"sync"

	"github.com/gorilla/websocket"
)

type PowerUp int32

type Action string

const (
	ChangeReady   = "ready"
	ChangeSetting = "config"
	PlayerCount   = "players"
	StatusChanged = "status"
	StatusRequest = "status_req"
	Submit        = "submit"
	SubmitFailed  = "submit_failed"
	UsePowerUp    = "use"
)

type ReadyState string

const (
	Ready   = "ready"
	Active  = "active"
	Waiting = "waiting"
)

type GameSettings struct {
	Language string
	Problem  string
}

type GameStatus string

const (
	Pending = "pending"
	Running = "active"
	End     = "end"
)

type GameState struct {
	sync.Mutex
	Status   GameStatus
	Settings GameSettings
}

func (gs *GameState) SetStatus(s GameStatus) {
	gs.Status = s
}

func (gs *GameState) SafeSetStatus(s GameStatus) {
	gs.Lock()
	defer gs.Unlock()

	gs.SetStatus(s)
}

func (gs *GameState) SetSettings(s GameSettings) {
	gs.Settings = s
}

func (gs *GameState) SafeSetSettings(s GameSettings) {
	gs.Lock()
	defer gs.Unlock()

	gs.SetSettings(s)
}

type PlayerState struct {
	Status      ReadyState
	Socket      *websocket.Conn
	SendMessage func(action Action, data interface{})
}
