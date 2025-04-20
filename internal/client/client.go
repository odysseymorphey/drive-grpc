package client

import (
	procs "drive-grpc/internal/client/grpc"
	"drive-grpc/internal/client/models"
	fs "drive-grpc/internal/pb"
	"drive-grpc/internal/storage/csv_storage"
	"drive-grpc/internal/utils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

type Client struct {
	Client     fs.FileServiceClient
	CC         *grpc.ClientConn
	csvStorage *csv_storage.CSVStorage
	cfg        *models.Config
}

func New(cfg *models.Config, csvStorage *csv_storage.CSVStorage) *Client {
	cc, err := grpc.NewClient(*cfg.ServerAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		logrus.Fatalf("failed to open connection: %v", err)
	}

	client := fs.NewFileServiceClient(cc)

	return &Client{
		Client:     client,
		CC:         cc,
		cfg:        cfg,
		csvStorage: csvStorage,
	}
}

func (c *Client) Run() {
	switch *c.cfg.Action {
	case "upload":
		if *c.cfg.FilePath == "" {
			logrus.Fatal("Wrong file path")
		}
		if !utils.IsImage(*c.cfg.FilePath) {
			logrus.Fatalf("File must be image")
		}

		csvRec, err := procs.UploadFile(c.Client, *c.cfg.FilePath)
		if err != nil {
			logrus.Fatalf("Upload failed: %v", err)
		}
		err = c.csvStorage.Save(csvRec)
		if err != nil {
			logrus.Fatalf("Failed to save CSV record: %v", err)
		}
	case "download":
		if *c.cfg.Filename == "" || *c.cfg.FileID == "" {
			logrus.Fatalf("Required file name and file id")
		}
		err := procs.DownloadFile(c.Client, *c.cfg.Filename, *c.cfg.FileID)
		if err != nil {
			log.Fatalf("Download failed: %v", err)
		}
	case "list":
		err := procs.FilesList(c.Client)
		if err != nil {
			log.Fatalf("Something went wrong: %v", err)
		}
	default:
		logrus.Fatal("Wrong action. Available actions: 'upload', 'download', 'list'")
	}
}

func (c *Client) Destroy() {
	c.CC.Close()
}
