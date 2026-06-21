package main

import (
	"log"
	"whiteboard/internal/server"
)

func main() {
	s := server.New()
	if err := s.Run(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
