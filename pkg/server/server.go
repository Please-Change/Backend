package server

import (
	"fmt"
	"net/http"

	"github.com/Please-Change/backend/pkg/types"
)

const PORT = 5174

var queue = make(chan types.GameState, 0)

func Serve() {
	fmt.Printf("Serving on %d\n", PORT)

	http.HandleFunc("/game", HandleStart)

	http.HandleFunc("/complete", HandleComplete)

	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
