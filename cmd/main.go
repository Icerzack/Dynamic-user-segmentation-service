package main

import (
	"avito-backend-internship/internal/server"
)

func main() {
	serverConfig := &server.Config{
		Port: "3001",
	}

	s := server.NewServer(serverConfig)

	s.Run()
}
