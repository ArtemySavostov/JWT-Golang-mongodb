package repository

import (
	"context"

	"github.com/ArtemySavostov/JWT-Golang-mongodb/internal/entity"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository interface {
	Get(ctx context.Context, id primitive.ObjectID) (entity.User, error)
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, user entity.User) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	GetByUsername(ctx context.Context, username string) (entity.User, error)
	GetAll(ctx context.Context) ([]entity.User, error)
}
