package handlers

import (
	"net/http"

	"JWT/internal/usecase"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userUC usecase.UserUseCase
}

func NewUserHandler(userUC usecase.UserUseCase) *UserHandler {
	return &UserHandler{userUC: userUC}
}

func (h *UserHandler) GetUser(c *gin.Context) {
	username := c.Param("username")

	user, err := h.userUC.GetUser(username)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
