package auth

import (
	"context"

	"log"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/ArtemySavostov/JWT-Golang-mongodb/internal/grpc"
	"github.com/ArtemySavostov/JWT-Golang-mongodb/internal/usecase"
)

type AuthGrpcServer struct {
	pb.UnimplementedAuthServiceServer
	authUC usecase.AuthUseCase
	userUC usecase.UserUseCase
}

func NewAuthGrpcServer(authUC usecase.AuthUseCase, userUC usecase.UserUseCase) *AuthGrpcServer {
	return &AuthGrpcServer{
		authUC: authUC,
		userUC: userUC,
	}
}

func (s *AuthGrpcServer) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

	if req.Username == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username and password are required")
	}

	_, err := s.userUC.CreateUser(req.Username, req.Email, req.Password)
	if err != nil {
		log.Printf("Error registering user: %v", err)

		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &pb.RegisterResponse{Success: true, Message: "User registered successfully"}, nil
}

func (s *AuthGrpcServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Username == "" || req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "username and password are required")
	}

	user, err := s.userUC.GetUser(req.Username)
	if err != nil {
		log.Printf("Error getting user: %v", err)
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	if !s.authUC.CheckPasswordHash(req.Password, user.Password) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}

	token, err := s.authUC.GenerateToken(user.Username, user.ID.Hex())
	if err != nil {
		log.Printf("Error generating token: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to generate token: %v", err)
	}

	return &pb.LoginResponse{Success: true, Token: token, Message: "Login successful"}, nil
}

func (s *AuthGrpcServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {

	if req.Token == "" {
		return nil, status.Errorf(codes.InvalidArgument, "token is required")
	}

	userID, err := s.authUC.ValidateToken(req.Token)
	if err != nil {
		log.Printf("Error validating token: %v", err)

		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}

	return &pb.ValidateTokenResponse{IsValid: true, UserID: userID}, nil
}
