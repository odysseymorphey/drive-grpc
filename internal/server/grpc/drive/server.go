package drive

import (
	files "drive-grpc/internal/pb"
	"drive-grpc/internal/server/repository"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Server struct {
	files.UnimplementedFileServiceServer
	repo         repository.Repository
	logger       *logrus.Logger
	filesStorage string
}

func New(repo repository.Repository) *Server {
	return &Server{
		repo:         repo,
		filesStorage: "/app/storage/",
		logger:       logrus.New(),
	}
}

func (s *Server) Run() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(UnaryInterceptor),
		grpc.StreamInterceptor(StreamInterceptor))

	files.RegisterFileServiceServer(grpcServer, s)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-shutdown
		log.Println("Starting graceful shutdown...")
		grpcServer.GracefulStop()
	}()

	s.logger.Info("gRPC server started at port 50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
