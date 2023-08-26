package main

import (
	"avito-backend-internship/internal/app"
)

func main() {
	serverConfig := &app.Config{
		Port: "3001",
	}

	s := app.NewServer(serverConfig)

	s.Run()
}
