package main

import (
	"log"
	"net/http"

	"github.com/ncalibey/learn-go-with-tests/http-server/server"
)

func main() {
	server := server.NewPlayerServer(server.NewInMemoryPlayerStore())

	if err := http.ListenAndServe(":5000", server); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
