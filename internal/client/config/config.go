package config

import (
	"drive-grpc/internal/client/models"
	"flag"
)

func New() *models.Config {
	cfg := &models.Config{}

	cfg.ServerAddr = flag.String("server", "localhost:50051", "server address")
	cfg.Action = flag.String("action", "", "upload|download|list")
	cfg.FilePath = flag.String("file", "", "path to file")
	cfg.OutputDir = flag.String("output", ".", "output directory")
	cfg.Filename = flag.String("filename", "", "filename for download")
	cfg.FileID = flag.String("id", "", "file id for download")
	flag.Parse()

	return cfg
}
