package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/try-go", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, world! %s", r.URL.Path)
	})

	http.ListenAndServe(":80", nil)
}
