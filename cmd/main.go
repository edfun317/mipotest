package main

import (
	"github.com/gin-gonic/gin"

	"github.com/edfun317/mipotest/internal/api/handler"
	"github.com/edfun317/mipotest/internal/api/router"
	"github.com/edfun317/mipotest/internal/infrastructure/db"
	"github.com/edfun317/mipotest/internal/infrastructure/persistence"
	"github.com/edfun317/mipotest/internal/usecase"
)

func main() {

	db := db.InitDB()
	//db.AutoMigrate(&entity.Account{}, &entity.TransactionLog{})

	// 初始化 Repository
	accRepo := persistence.NewAccountGormRepository(db)
	txRepo := persistence.NewTransactionGormRepository(db)

	//  in Usecase
	uc := &usecase.AccountUsecase{
		DB:              db,
		AccountRepo:     accRepo,
		TransactionRepo: txRepo,
	}

	// initialize Handler
	h := &handler.AccountHandler{Usecase: uc}

	// register routes
	r := gin.Default()
	router.RegisterRoutes(r, h)

	// run server
	r.Run(":8080")
}
