package usecase

import (
	"JWT/internal/entity"
	"JWT/internal/repository"
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type UserUseCase interface {
	GetUser(username string) (entity.User, error)
	CreateUser(username, email, password string) (entity.User, error)
	UpdateUser(user entity.User) error
	DeleteUser(id primitive.ObjectID) error
	GetAllUsers() ([]entity.User, error)
}

type userUseCase struct {
	userRepo repository.UserRepository
}

func NewUserUseCase(userRepo repository.UserRepository) UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (uc *userUseCase) GetUser(username string) (entity.User, error) {
	return uc.userRepo.GetByUsername(context.Background(), username)
}

func (uc *userUseCase) CreateUser(username, email, password string) (entity.User, error) {
	if username == "" || email == "" || password == "" {
		return entity.User{}, errors.New("username, email, and password are required")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return entity.User{}, fmt.Errorf("could not hash password: %w", err)
	}

	user := &entity.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	err = uc.userRepo.Create(context.Background(), user)
	if err != nil {
		return entity.User{}, fmt.Errorf("could not create user: %w", err)
	}

	return *user, nil
}

func (uc *userUseCase) UpdateUser(user entity.User) error {
	return uc.userRepo.Update(context.Background(), user)
}

func (uc *userUseCase) DeleteUser(id primitive.ObjectID) error {
	return uc.userRepo.Delete(context.Background(), id)
}

func (uc *userUseCase) GetAllUsers() ([]entity.User, error) {
	return uc.userRepo.GetAll(context.Background())
}
