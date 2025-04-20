package csv_storage

import (
	"encoding/csv"
	"github.com/odysseymorphey/drive-grpc/internal/server/models"
	"os"
)

type CSVStorage struct {
	FilePath string
}

func New(filePath string) *CSVStorage {
	return &CSVStorage{
		FilePath: filePath,
	}
}

func (s *CSVStorage) Save(record *models.CSVRecord) error {
	file, err := os.OpenFile(s.FilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	if fileInfo, err := file.Stat(); fileInfo.Size() == 0 && err != nil {
		header := []string{"fileId", "filename"}
		if err := writer.Write(header); err != nil {
			return err
		}
	}

	data := []string{
		record.FileID,
		record.FileName,
	}

	return writer.Write(data)
}
