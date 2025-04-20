package drive

import (
	"context"
	"github.com/google/uuid"
	fs "github.com/odysseymorphey/drive-grpc/internal/pb"
	"github.com/odysseymorphey/drive-grpc/internal/server/models"
	"github.com/odysseymorphey/drive-grpc/internal/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
	"io"
	"os"
	"path/filepath"
	"time"
)

const chunkSize = 32 * 1024

func (s *Server) UploadFile(stream fs.FileService_UploadFileServer) error {
	ctx := stream.Context()

	req, err := stream.Recv()
	if err != nil {
		s.logger.Errorf("cannot receive file info: %v", err)
		return status.Error(codes.Internal, "stream error")
	}

	filename := utils.SanitizeFilename(req.GetFilename())
	if filename == "" {
		s.logger.Errorf("invalid file name. get: %v", filename)
		return status.Error(codes.InvalidArgument, "invalid file name")
	}

	folder := uuid.New().String()
	dirPath := filepath.Join(s.filesStorage, folder)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		s.logger.Errorf("failed to create directory: %v", err)
		return status.Error(codes.Internal, "storage error")
	}

	filePath := filepath.Join(dirPath, filename)
	file, err := os.Create(filePath)
	if err != nil {
		s.logger.Errorf("failed to create file: %v", err)
		return status.Error(codes.Internal, "storage error")
	}
	defer file.Close()

	var size uint32

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			s.logger.Errorf("failed to receive chunk: %v", err)
			return status.Error(codes.Internal, "chunk receive error")
		}

		chunk := req.GetChunk()
		n, err := file.Write(chunk)
		size += uint32(n)
	}

	fileInfo := &models.FileInfo{
		ID:               uuid.New().String(),
		Filename:         filename,
		FilePath:         filePath,
		CreationDate:     time.Now(),
		ModificationDate: time.Now(),
	}

	err = s.repo.WriteFileInfo(ctx, fileInfo)
	if err != nil {
		return err
	}

	return stream.SendAndClose(&fs.UploadFileResponse{
		Filename: filename,
		Id:       fileInfo.ID,
		Size:     size,
	})
}

func (s *Server) FilesList(
	ctx context.Context,
	req *fs.FilesListRequest) (*fs.FilesListResponse, error) {
	fInfoList, err := s.repo.GetFilesInfoList(ctx)
	if err != nil {
		s.logger.Errorf("failed to read data from db: %v", err)
		return nil, status.Error(codes.Internal, "read failed")
	}

	var files []*fs.FileInfo

	for _, fInfo := range fInfoList {
		files = append(files, &fs.FileInfo{
			Id:               fInfo.ID,
			Filename:         fInfo.Filename,
			CreationDate:     timestamppb.New(fInfo.CreationDate),
			ModificationDate: timestamppb.New(fInfo.ModificationDate),
		})
	}

	return &fs.FilesListResponse{
		Files: files,
	}, nil
}

func (s *Server) DownloadFile(
	req *fs.DownloadFileRequest,
	stream fs.FileService_DownloadFileServer) error {
	ctx := stream.Context()

	filename := req.GetFilename()
	fileID := req.GetId()
	fInfo, err := s.repo.GetFileInfoByName(ctx, filename, fileID)
	if err != nil {
		s.logger.Errorf("failed to get file info: %v", err)
		return status.Error(codes.Internal, "download failed")
	}

	file, err := os.Open(fInfo.FilePath)
	if err != nil {
		s.logger.Errorf("failed to open file: %v", err)
		return status.Error(codes.Internal, "storage error")
	}
	defer file.Close()

	err = stream.Send(&fs.DownloadFileResponse{
		Filename: filename,
	})
	if err != nil {
		s.logger.Errorf("failed to send first message")
		return status.Error(codes.Internal, "download failed")
	}

	buff := make([]byte, chunkSize)
	for {
		select {
		case <-ctx.Done():
			return status.FromContextError(ctx.Err()).Err()
		default:
			n, err := file.Read(buff)
			if err == io.EOF {
				return nil
			}
			if err != nil {
				s.logger.Errorf("read error: %v", err)
				return status.Error(codes.Internal, "read failed")
			}

			chunk := make([]byte, n)
			copy(chunk, buff[:n])

			err = stream.Send(&fs.DownloadFileResponse{
				Chunk: chunk,
			})
			if err != nil {
				s.logger.Errorf("failed to send data: %v", err)
				return status.Error(codes.Internal, "send failed")
			}
		}
	}
}
