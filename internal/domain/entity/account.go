package entity

import "time"

type Account struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:255;not null"`
	Balance   int64  `gorm:"not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type TransactionLog struct {
	ID            uint `gorm:"primaryKey"`
	FromAccountID *uint
	ToAccountID   *uint
	Amount        int64 // positive only
	CreatedAt     time.Time
	Description   string
}
