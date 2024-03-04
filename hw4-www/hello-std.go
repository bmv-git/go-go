package main

import (
	"fmt"
	"log"
	"net/http"
)

func hello(w http.ResponseWriter, _ *http.Request) {
	_, err := fmt.Fprintf(w, "Hello World!")
	if err != nil {
		return
	}
}
func main() {
	port := ":8080"
	http.HandleFunc("/", hello)
	log.Printf("Server is starting on port :%v", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
