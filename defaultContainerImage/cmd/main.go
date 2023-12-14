package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w,
			"Default Container Image\nResource Group: %s\nApplication Name: %s",
			os.Getenv("RESOURCE_GROUP"),
			os.Getenv("APPLICATION_NAME"),
		)
	})

	http.ListenAndServe(":8080", nil)
}
