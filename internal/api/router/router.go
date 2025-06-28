package router

import (
	"github.com/edfun317/mipotest/internal/api/handler"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, h *handler.AccountHandler) {

	r.POST("/accounts", h.CreateAccount)
	r.POST("/accounts/:id/deposit", h.Deposit)
	r.POST("/accounts/:id/withdraw", h.Withdraw)
	r.POST("/accounts/transfer", h.Transfer)
	r.GET("/accounts/:id", h.GetAccount)
	r.GET("/accounts/:id/logs", h.GetTransactionLogs)
}
