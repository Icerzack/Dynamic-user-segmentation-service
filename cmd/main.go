package main

import (
	"avito-backend-internship/internal/app"
	"avito-backend-internship/internal/pkg/db"
	database "avito-backend-internship/internal/pkg/service/db"
	"avito-backend-internship/internal/pkg/service/history"
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

	postgresService := database.NewPostgresService()
	csvHistoryService := history.NewCSVHistoryService()

	s := app.NewServer(
		ctx,
		config,
		postgresDB,
		postgresService,
		csvHistoryService,
	)

	s.Run()
}
