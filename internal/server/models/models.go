package models

import (
	"time"
)

type FileInfo struct {
	ID               string
	Filename         string
	FilePath         string
	CreationDate     time.Time
	ModificationDate time.Time
}

type CSVRecord struct {
	FileID   string
	FileName string
}
