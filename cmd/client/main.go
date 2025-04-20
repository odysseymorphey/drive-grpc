package main

import (
	"drive-grpc/internal/client"
	"drive-grpc/internal/client/config"
	"drive-grpc/internal/storage/csv_storage"
)

func main() {
	cfg := config.New()
	csvStorage := csv_storage.New("./files_list.csv")
	cl := client.New(cfg, csvStorage)
	cl.Run()
	cl.Destroy()
}
