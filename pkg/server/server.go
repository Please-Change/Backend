package server

import (
	"fmt"
	"net/http"
)

const PORT = 5174

func Serve() {
	fmt.Printf("Serving on %d\n", PORT)

	http.HandleFunc("/game", HandleStart)

	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
