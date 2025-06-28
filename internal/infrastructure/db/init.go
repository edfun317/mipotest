package db

import (
	"sync"

	"github.com/edfun317/mipotest/internal/domain/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	once       sync.Once
	dbInstance *gorm.DB
)

func InitDB() *gorm.DB {

	db, err := gorm.Open(sqlite.Open("bank.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	return db
}

func AutoMigrate(db *gorm.DB) {
	err := db.AutoMigrate(&entity.Account{}, &entity.TransactionLog{})
	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {

	//singleton pattern to ensure only one instance of the database connection
	once.Do(func() {
		dbInstance = InitDB()
	})

	return dbInstance
}
