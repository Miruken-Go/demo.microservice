package main

import (
    "fmt"
	"os"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Default Container Image\nENV: %s\n", os.Getenv("ENV"))
    })

    http.ListenAndServe(":8080", nil)
}