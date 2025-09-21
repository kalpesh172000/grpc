package main

import (
    "context"
    "log"

    pb "github.com/kalpesh172000/grpc/proto"
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
)

func main() {
    // Connect to server
    conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatalf("Failed to connect: %v", err)
    }
    defer conn.Close()

    client := pb.NewAuthServiceClient(conn)

    // Register a user
    registerResp, err := client.Register(context.Background(), &pb.RegisterRequest{
        Email:    "test@example.com",
        Password: "password123",
        Name:     "Test User",
    })
    if err != nil {
        log.Fatalf("Register failed: %v", err)
    }
    log.Printf("Register: %+v", registerResp)

    // Login
    loginResp, err := client.Login(context.Background(), &pb.LoginRequest{
        Email:    "test@example.com",
        Password: "password123",
    })
    if err != nil {
        log.Fatalf("Login failed: %v", err)
    }
    log.Printf("Login: %+v", loginResp)

    // Get profile
    profileResp, err := client.GetProfile(context.Background(), &pb.GetProfileRequest{
        Token: loginResp.Token,
    })
    if err != nil {
        log.Fatalf("GetProfile failed: %v", err)
    }
    log.Printf("Profile: %+v", profileResp)
}
