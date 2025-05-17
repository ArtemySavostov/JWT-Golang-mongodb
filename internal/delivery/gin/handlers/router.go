package handlers

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(authHandler *AuthHandler, userHandler *UserHandler) *gin.Engine {
	r := gin.Default()

	r.POST("/register", authHandler.Register)
	r.POST("/login", authHandler.Login)

	r.GET("/users/:username", AuthMiddleware(), userHandler.GetUser)

	return r
}
