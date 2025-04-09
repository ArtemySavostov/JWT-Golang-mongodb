package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"JWT/internal/app"
	"JWT/internal/delivery/gin"
	"JWT/internal/delivery/gin/handlers"
	"JWT/internal/usecase"
	"JWT/pkg/database"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI environment variable not set")
	}

	client, err := database.ConnectMongoDB(context.Background(), mongoURI)
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatalf("Error disconnecting from MongoDB: %v", err)
		}
	}()

	userRepo := app.NewMongoDBUserRepository(client.Database("forum_db").Collection("users"))
	userUC := usecase.NewUserUseCase(userRepo)
	authUC := usecase.NewAuthUseCase(userRepo)
	userHandler := handlers.NewUserHandler(userUC)
	authHandler := handlers.NewAuthHandler(authUC)

	router := gin.SetupRouter(authHandler, userHandler)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})
	handler := c.Handler(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	fmt.Printf("Server is running on port %s...\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
