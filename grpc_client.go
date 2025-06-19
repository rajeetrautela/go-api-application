package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	pb "go-jwt-api/internal/fileupload"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func UploadFileToGRPCServer(filePath string) (string, error) {

	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return "", fmt.Errorf("failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	client := pb.NewFileUploadServiceClient(conn)

	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %v", err)
	}

	req := &pb.FileRequest{
		Filename: filepath.Base(filePath),
		Content:  content,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	res, err := client.UploadFile(ctx, req)
	if err != nil {
		return "", fmt.Errorf("upload failed: %v", err)
	}

	return res.Message, nil
}
