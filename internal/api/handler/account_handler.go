package handler

import (
	"net/http"

	"github.com/edfun317/mipotest/internal/usecase"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	Usecase *usecase.AccountUsecase
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {

	var req struct {
		Name    string `json:"name"`
		Balance int64  `json:"balance"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	acc, err := h.Usecase.CreateAccount(req.Name, req.Balance)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, acc)
}
