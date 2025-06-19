package config

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"

	pb "go-jwt-api/internal/fileupload"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedFileUploadServiceServer
}

func (s *server) UploadFile(ctx context.Context, req *pb.FileRequest) (*pb.FileResponse, error) {
	os.MkdirAll("uploads", os.ModePerm)
	dstPath := filepath.Join("uploads", req.Filename)
	err := os.WriteFile(dstPath, req.Content, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to save file: %v", err)
	}
	return &pb.FileResponse{Message: "File uploaded successfully"}, nil
}

func StartGRPCServer() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}
	grpcServer := grpc.NewServer()
	pb.RegisterFileUploadServiceServer(grpcServer, &server{})
	fmt.Println("gRPC server listening on :50051")
	grpcServer.Serve(lis)
}
