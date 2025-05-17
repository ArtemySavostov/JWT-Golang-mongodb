package app

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/ArtemySavostov/JWT-Golang-mongodb/internal/entity"
	"github.com/ArtemySavostov/JWT-Golang-mongodb/internal/repository"

	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoDBUserRepository struct {
	collection *mongo.Collection
}

func NewMongoDBUserRepository(collection *mongo.Collection) repository.UserRepository {
	return &MongoDBUserRepository{collection: collection}
}

func (r *MongoDBUserRepository) Get(ctx context.Context, id primitive.ObjectID) (entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, fmt.Errorf("user not found")
		}
		log.Printf("failed to FindOne on db: %v", err)
		return entity.User{}, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *MongoDBUserRepository) Create(ctx context.Context, user *entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.collection.InsertOne(ctx, user)

	if err != nil {
		log.Printf("failed to InsertOne on db: %v", err)
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *MongoDBUserRepository) Update(ctx context.Context, user entity.User) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": user.ID}
	update := bson.M{
		"$set": bson.M{
			"username": user.Username,
			"email":    user.Email,
			"password": user.Password,
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("failed to UpdateOne on db: %v", err)
		return fmt.Errorf("failed to update user: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *MongoDBUserRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": id}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("failed to DeleteOne on db: %v", err)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *MongoDBUserRepository) GetByUsername(ctx context.Context, username string) (entity.User, error) {
	var user entity.User
	err := r.collection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return entity.User{}, fmt.Errorf("user not found")
		}
		return entity.User{}, fmt.Errorf("failed to get user by username: %w", err)
	}

	return user, nil
}

func (r *MongoDBUserRepository) GetAll(ctx context.Context) ([]entity.User, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("failed to Find on db: %v", err)
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}
	defer cursor.Close(ctx)

	var users []entity.User

	for cursor.Next(ctx) {
		var user entity.User

		if err := cursor.Decode(&user); err != nil {
			log.Printf("failed to decode user: %v", err)
			return nil, fmt.Errorf("failed to decode user: %w", err)
		}

		users = append(users, user)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("cursor error: %v", err)
		return nil, fmt.Errorf("cursor error: %w", err)
	}

	return users, nil
}
