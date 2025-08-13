package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/ryszhio/tasktracker/internal/server"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("No .env file found")
	}
	server := server.NewServer()

	log.Fatal(server.ListenAndServe())
}
