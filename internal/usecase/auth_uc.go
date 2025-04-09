package usecase

import (
	"JWT/internal/entity"
	"JWT/internal/repository"
	"JWT/pkg/auth"
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type AuthUseCase interface {
	Register(ctx context.Context, username, email, password string) (string, error)
	Login(ctx context.Context, username, password string) (string, error)
	ValidateToken(token string) (string, error)
	GenerateToken(username string, userID string) (string, error)
	CheckPasswordHash(password string, hashedPassword string) bool
}

type authUseCase struct {
	userRepo repository.UserRepository
}

func NewAuthUseCase(userRepo repository.UserRepository) AuthUseCase {
	return &authUseCase{userRepo: userRepo}
}

func (uc *authUseCase) Register(ctx context.Context, username, email, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("could not hash password: %w", err)
	}

	user := &entity.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	err = uc.userRepo.Create(ctx, user)
	if err != nil {
		return "", fmt.Errorf("could not create user: %w", err)
	}

	token, err := uc.GenerateToken(username, user.ID.Hex())
	if err != nil {
		return "", fmt.Errorf("could not generate token: %w", err)
	}

	return token, nil
}

func (uc *authUseCase) Login(ctx context.Context, username, password string) (string, error) {
	user, err := uc.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", fmt.Errorf("could not get user: %w", err)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := uc.GenerateToken(username, user.ID.Hex())
	if err != nil {
		return "", fmt.Errorf("could not generate token: %w", err)
	}

	return token, nil
}

func (uc *authUseCase) ValidateToken(token string) (string, error) {
	_, userID, err := auth.ValidateToken(token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	return userID, nil
}

func (uc *authUseCase) GenerateToken(username string, userID string) (string, error) {
	token, err := auth.GenerateToken(username, userID)
	if err != nil {
		return "", fmt.Errorf("could not generate token: %w", err)
	}

	return token, nil
}

func (uc *authUseCase) CheckPasswordHash(password string, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
