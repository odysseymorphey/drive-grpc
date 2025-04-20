package repository

import (
	"context"
	"github.com/odysseymorphey/drive-grpc/internal/server/models"
)

type Repository interface {
	GetFileInfoByName(ctx context.Context, filename string, fileId string) (*models.FileInfo, error)
	GetFilesInfoList(ctx context.Context) ([]models.FileInfo, error)
	WriteFileInfo(ctx context.Context, fileInfo *models.FileInfo) error
}
