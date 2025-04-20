package main

import (
	"drive-grpc/internal/server/grpc/drive"
	"drive-grpc/internal/storage/postgres"
	"fmt"
	"log"
	"os"
)

func main() {
	db, err := postgres.New(fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
	),
	)
	if err != nil {
		log.Fatalf("Can't open database: %v", err)
	}

	app := drive.New(db)
	app.Run()
}
