package handlers

import (
	"net/http"
	"task-management-backend/internal/domain/entity"
	"task-management-backend/internal/usecase/auth"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authUC *auth.AuthUseCase
}

func NewAuthHandler(authUC *auth.AuthUseCase) *AuthHandler {
	return &AuthHandler{
		authUC: authUC,
	}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req entity.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authUC.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
