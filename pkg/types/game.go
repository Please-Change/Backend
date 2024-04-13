package types

import "github.com/gorilla/websocket"

type PowerUp int32

type Action string

const (
	UsePowerUp    string = "use"
	ChangeReady          = "ready"
	ChangeSetting        = "config"
	Submit               = "submit"
	StatusRequest        = "status_req"
	PlayerCount          = "players"
)

type ReadyState string

const (
	Ready   string = "ready"
	Active         = "active"
	Waiting        = "waiting"
)

type GameSettings struct{}

type GameState struct {
	Ready    ReadyState
	Settings GameSettings
	Socket   *websocket.Conn
}
