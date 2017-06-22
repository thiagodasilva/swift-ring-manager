package main

import (
	"log"
	"net/http"

	"github.com/thiagodasilva/swift-ring-manager/ringmanager"
)

func main() {

	router := ringmanager.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
