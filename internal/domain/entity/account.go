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
	ID            uint  `gorm:"primaryKey"`
	FromAccountID *uint `gorm:"index:idx_from_account_time"`
	ToAccountID   *uint `gorm:"index:idx_to_account_time"`
	Amount        int64
	CreatedAt     time.Time `gorm:"index:idx_from_account_time;index:idx_to_account_time"`
	Description   string
}
