package main

import (
	"log"
	"net/http"

	"github.com/thiagodasilva/swift-ring-manager/pkg/ringmanager"
)

func main() {

	router := ringmanager.NewRouter()
	log.Fatal(http.ListenAndServe(":8080", router))
}
