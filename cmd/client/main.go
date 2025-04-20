package main

import (
	"github.com/odysseymorphey/drive-grpc/internal/client"
	"github.com/odysseymorphey/drive-grpc/internal/client/config"
	"github.com/odysseymorphey/drive-grpc/internal/storage/csv_storage"
)

func main() {
	cfg := config.New()
	csvStorage := csv_storage.New("./files_list.csv")
	cl := client.New(cfg, csvStorage)
	cl.Run()
	cl.Destroy()
}
