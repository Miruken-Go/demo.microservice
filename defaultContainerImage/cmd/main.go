package main

import (
    "fmt"
	"os"
    "net/http"
)

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Default Container Image\nResource Group: %s\n", os.Getenv("RESOURCE_GROUP"))
    })

    http.ListenAndServe(":8080", nil)
}