package main

import (
	"avito-backend-internship/internal/app"
	"avito-backend-internship/internal/pkg/db"
	"avito-backend-internship/internal/pkg/service"
	"context"
	"fmt"
)

func main() {
	config := &app.Config{
		Port: "3001",
	}

	ctx := context.Background()

	postgresDB, err := db.NewDB(ctx)
	if err != nil {
		fmt.Println("Error while setting up a database:", err.Error())
		return
	}

	postgresService := service.NewPostgresService()

	s := app.NewServer(
		ctx,
		config,
		postgresDB,
		postgresService,
	)

	s.Run()
}
