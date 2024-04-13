package types

type GameState struct {
	gameId       ObjectID
	targetTextId ObjectID
	currentText  chan string
	powerups     []PowerUp
	isPaused     bool
}

type PowerUp string

const (
	SpeechToText PowerUp = "speechtotext"
	TextToSpeech         = "texttospeech"
	SwapKeys             = "swapkeys"
	BanKey               = "bankey"
)

type Action string

const (
	KeyPress    string = "key"
	UsePowerUp         = "use"
	QuitGame           = "quit"
	TogglePause        = "stop"
)
