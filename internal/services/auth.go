package services

import (
	"context"
	"errors"
	"time"

	"github.com/kalpesh172000/grpc/internal/database"
	"github.com/kalpesh172000/grpc/internal/models"
	"github.com/kalpesh172000/grpc/internal/utils"
	pb "github.com/kalpesh172000/grpc/proto"
	"gorm.io/gorm"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	jwtSecret string
	logger    *utils.Logger
}

func NewAuthServer(jwtSecret string, logger *utils.Logger) *AuthServer {
	return &AuthServer{
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (s *AuthServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	start := time.Now()
	defer s.logger.LogRequest("gRPC", "Register", "internal", 200, time.Since(start))

	// Check if user already exists
	var existingUser models.User
	if err := database.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: "User already exists",
		}, nil
	}

	// Create new user
	user := models.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	// Hash password
	if err := user.HashPassword(); err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: "Failed to hash password",
		}, err
	}

	// Save to database
	if err := database.DB.Create(&user).Error; err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: "Failed to create user",
		}, err
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return &pb.RegisterResponse{
			Success: false,
			Message: "Failed to generate token",
		}, err
	}

	return &pb.RegisterResponse{
		Success: true,
		Message: "User registered successfully",
		Token:   token,
		User: &pb.User{
			Id:        uint32(user.ID),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *AuthServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	start := time.Now()
	defer s.logger.LogRequest("gRPC", "Login", "internal", 200, time.Since(start))

	// Find user by email
	var user models.User
	if err := database.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.LoginResponse{
				Success: false,
				Message: "Invalid credentials",
			}, nil
		}
		return &pb.LoginResponse{
			Success: false,
			Message: "Database error",
		}, err
	}

	// Check password
	if !user.CheckPassword(req.Password) {
		return &pb.LoginResponse{
			Success: false,
			Message: "Invalid credentials",
		}, nil
	}

	// Generate JWT token
	token, err := utils.GenerateToken(user.ID, user.Email, s.jwtSecret)
	if err != nil {
		return &pb.LoginResponse{
			Success: false,
			Message: "Failed to generate token",
		}, err
	}

	return &pb.LoginResponse{
		Success: true,
		Message: "Login successful",
		Token:   token,
		User: &pb.User{
			Id:        uint32(user.ID),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}

func (s *AuthServer) GetProfile(ctx context.Context, req *pb.GetProfileRequest) (*pb.GetProfileResponse, error) {
	start := time.Now()
	defer s.logger.LogRequest("gRPC", "GetProfile", "internal", 200, time.Since(start))

	// Validate JWT token
	claims, err := utils.ValidateToken(req.Token, s.jwtSecret)
	if err != nil {
		return &pb.GetProfileResponse{
			Success: false,
			Message: "Invalid token",
		}, nil
	}

	// Find user by ID
	var user models.User
	if err := database.DB.First(&user, claims.UserID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &pb.GetProfileResponse{
				Success: false,
				Message: "User not found",
			}, nil
		}
		return &pb.GetProfileResponse{
			Success: false,
			Message: "Database error",
		}, err
	}

	return &pb.GetProfileResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		User: &pb.User{
			Id:        uint32(user.ID),
			Email:     user.Email,
			Name:      user.Name,
			CreatedAt: user.CreatedAt.Format(time.RFC3339),
			UpdatedAt: user.UpdatedAt.Format(time.RFC3339),
		},
	}, nil
}
