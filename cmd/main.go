package main

import (
	"log"
	"net/http"
	"os"

	"github.com/ncalibey/learn-go-with-tests/http-server/internal/server"
)

const dbFileName = "game.db.json"

func main() {
	db, err := os.OpenFile(dbFileName, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalf("problem opening %s %v", dbFileName, err)
	}

	store, err := server.NewFilesystemPlayerStore(db)
	if err != nil {
		log.Fatalf("problem creating file system player store, %v ", err)
	}

	s := server.NewPlayerServer(store)

	if err := http.ListenAndServe(":5000", s); err != nil {
		log.Fatalf("could not listen on port 5000 %v", err)
	}
}
