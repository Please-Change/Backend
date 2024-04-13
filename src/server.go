package backend

import (
	"fmt"
	"net/http"
)

const PORT = 3000

func serve() {
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	fmt.Printf("Serving on %d", PORT)
}
