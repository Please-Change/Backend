package types

type GameState struct {
	targetTextId int32
	currentText  chan string
	powerups     []PowerUp
}

type PowerUp int32

const (
	Undefined PowerUp = iota
	SpeechToText
	TextToSpeech
	SwapKeys
	BanKey
)

type Action string

const (
	KeyPress   string = "key"
	UsePowerUp        = "use"
	QuitGame          = "quit"
	Pause             = "stop"
)
