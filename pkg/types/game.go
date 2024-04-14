package types

import (
	"bytes"
	"sync"

	"github.com/bytedance/sonic"
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

type Language string

const (
	C          Language = "c"
	Cpp                 = "cpp"
	Go                  = "go"
	Python              = "python"
	JavaScript          = "javascript"
	Rust                = "rust"
)

type GameSettings struct {
	Language Language
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
	sync.Mutex
	Status ReadyState
	Socket *websocket.Conn
}

func (ps *PlayerState) SetStatus(s ReadyState) {
	ps.Status = s
}

func (ps *PlayerState) SafeSetStatus(s ReadyState) {
	ps.Lock()
	defer ps.Unlock()

	ps.SetStatus(s)
}

func (ps *PlayerState) SendMessage(mt int, action Action, data interface{}) {
	var msg = map[string]interface{}{
		"action": action,
		"data":   data,
	}

	var w = bytes.NewBuffer(nil)
	var enc = sonic.ConfigDefault.NewEncoder(w)
	enc.Encode(msg)
	ps.Socket.WriteMessage(mt, w.Bytes())
}
