package types

import "github.com/gorilla/websocket"

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
	Status   GameStatus
	Settings GameSettings
}

type PlayerState struct {
	Status ReadyState
	Socket *websocket.Conn
}
