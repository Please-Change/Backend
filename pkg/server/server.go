package server

import (
	"fmt"
	"net/http"
	"sync"
)

const PORT = 5174

var QueueGroup = new(sync.WaitGroup)
var CurrentlyInGame = make(chan bool)
var ReadyPlayerCount = make(chan int64)

func Serve() {
	QueueGroup.Add(MAX_QUEUED_PLAYERS)
	fmt.Printf("Serving on %d\n", PORT)

	http.HandleFunc("/game", HandleStart)

	http.HandleFunc("/complete", HandleComplete)

	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
