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

func Serve() {
	QueueGroup.Add(MAX_QUEUED_PLAYERS)
	fmt.Printf("Serving on %d\n", PORT)

	http.HandleFunc("/game", HandleStart)

	http.HandleFunc("/complete", HandleComplete)

	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
