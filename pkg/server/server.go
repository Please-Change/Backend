package server

import (
	"fmt"
	"net/http"
)

const PORT = 3000

func Serve() {
	fmt.Printf("Serving on %d", PORT)

	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)
}
