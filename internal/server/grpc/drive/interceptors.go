package drive

import (
	"context"
	"google.golang.org/grpc"
)

var (
	updownLimiter = make(chan struct{}, 10)
	viewLimiter   = make(chan struct{}, 100)
)

func UnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	switch info.FullMethod {
	case "/fs.FileService/FilesList":
		viewLimiter <- struct{}{}
		defer func() { <-viewLimiter }()
	case "/fs.FileServer/DownloadFile":
		updownLimiter <- struct{}{}
		defer func() { <-updownLimiter }()
	}

	return handler(ctx, req)
}

func StreamInterceptor(
	srv interface{},
	ss grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	switch info.FullMethod {
	case "/fs.FileService/UploadFile":
		updownLimiter <- struct{}{}
		defer func() { <-updownLimiter }()
	}

	return handler(srv, ss)
}
