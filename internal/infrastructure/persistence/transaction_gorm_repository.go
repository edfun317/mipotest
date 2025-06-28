package persistence

import (
	"github.com/edfun317/mipotest/internal/domain/entity"
	"github.com/edfun317/mipotest/internal/domain/repository"

	"gorm.io/gorm"
)

type transactionGormRepository struct{ db *gorm.DB }

func NewTransactionGormRepository(db *gorm.DB) repository.TransactionRepository {
	return &transactionGormRepository{db}
}

func (r *transactionGormRepository) Create(t *entity.TransactionLog) error {
	return r.db.Create(t).Error
}
func (r *transactionGormRepository) FindByAccount(accountID uint) ([]*entity.TransactionLog, error) {
	var logs []*entity.TransactionLog
	if err := r.db.
		Where("from_account_id = ? OR to_account_id = ?", accountID, accountID).
		Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}
