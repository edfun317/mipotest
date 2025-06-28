package repository

import "github.com/edfun317/mipotest/internal/domain/entity"

type TransactionRepository interface {
	Create(log *entity.TransactionLog) error
	FindByAccount(accountID uint) ([]*entity.TransactionLog, error)
}
