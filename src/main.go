package main

import (
	"api-note/src/server"
	"log"

	"github.com/gorilla/sessions"
)

func main() {
	log.SetFlags(log.Lshortfile)

	sessionStore := sessions.NewCookieStore([]byte("sessionKey"))
	server := server.NewServer(":8080", sessionStore)

	if err := server.Listen(); err != nil {
		log.Fatal(err)
	}
}
