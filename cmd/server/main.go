package main

import (
	"log"
	"net"

	"github.com/kalpesh172000/grpc/internal/config"
	"github.com/kalpesh172000/grpc/internal/database"
	"github.com/kalpesh172000/grpc/internal/services"
	"github.com/kalpesh172000/grpc/internal/utils"
	pb "github.com/kalpesh172000/grpc/proto"
	"google.golang.org/grpc"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()
	
	// Initialize logger
	logger := utils.NewLogger(cfg.LogFile)
	defer logger.Close()

	// Connect to database
	database.Connect(cfg.DatabaseURL)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register auth service
	authServer := services.NewAuthServer(cfg.JWTSecret, logger)
	pb.RegisterAuthServiceServer(grpcServer, authServer)

	// Start listening
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("gRPC server starting on port %s...", cfg.Port)
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
