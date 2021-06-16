package main

import (
	"api-note/src/server"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)

	server := server.NewServer(":8080")

	if err := server.Listen(); err != nil {
		log.Fatal(err)
	}
}
