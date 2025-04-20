package grpc

import (
	"context"
	fs "drive-grpc/internal/pb"
	"drive-grpc/internal/server/models"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	chunkSize = 32 * 1024

	timeOut = 30 * time.Second
)

func UploadFile(cl fs.FileServiceClient, filePath string) (*models.CSVRecord, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %s: %v", filePath, err)
	}
	defer file.Close()

	stream, err := cl.UploadFile(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to open stream: %v", err)
	}

	filename := filepath.Base(filePath)

	err = stream.Send(&fs.UploadFileRequest{
		Filename: filename,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send fisrt message: %v", err)
	}

	buff := make([]byte, chunkSize)
	for {
		n, err := file.Read(buff)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read error: %v", err)
		}

		chunk := make([]byte, n)
		copy(chunk, buff)

		err = stream.Send(&fs.UploadFileRequest{
			Chunk: chunk,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to send data message: %v", err)
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		return nil, fmt.Errorf("upload failed: %v", err)
	}

	csvRecord := &models.CSVRecord{
		FileID:   res.Id,
		FileName: res.Filename,
	}

	log.Printf("File uploaded. FileID: %v, Filename: %v. Size: %v", res.Id, res.Filename, res.Size)

	return csvRecord, nil
}

func FilesList(cl fs.FileServiceClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	res, err := cl.FilesList(ctx, &fs.FilesListRequest{})
	if err != nil {
		return fmt.Errorf("failed to get files list: %v", err)
	}

	for _, v := range res.Files {
		fmt.Printf("%s | %s | %s",
			v.Filename,
			v.CreationDate.AsTime().Format(time.RFC3339),
			v.ModificationDate.AsTime().Format(time.RFC3339),
		)
	}

	return nil
}

func DownloadFile(client fs.FileServiceClient, filename string, fileId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeOut)
	defer cancel()

	stream, err := client.DownloadFile(ctx, &fs.DownloadFileRequest{
		Filename: filename,
		Id:       fileId,
	})
	if err != nil {
		return fmt.Errorf("failed to open stream: %v", err)
	}

	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("download failed: %v", err)
		}

		if _, err = file.Write(res.Chunk); err != nil {
			return fmt.Errorf("failed to write chunk: %v", err)
		}
	}

	return nil
}
