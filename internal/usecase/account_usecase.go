package usecase

import (
	"errors"
	"time"

	"github.com/edfun317/mipotest/internal/domain/entity"
	"github.com/edfun317/mipotest/internal/domain/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AccountUsecase struct {
	DB              *gorm.DB
	AccountRepo     repository.AccountRepository
	TransactionRepo repository.TransactionRepository
}

func (u *AccountUsecase) CreateAccount(name string, balance int64) (*entity.Account, error) {
	if balance < 0 {
		return nil, errors.New("initial balance cannot be negative")
	}
	acc := &entity.Account{Name: name, Balance: balance}
	if err := u.AccountRepo.Create(acc); err != nil {
		return nil, err
	}
	return acc, nil
}

func (u *AccountUsecase) Deposit(accountID uint, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	return u.DB.Transaction(func(tx *gorm.DB) error {
		var acc entity.Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&acc, accountID).Error; err != nil {
			return err
		}
		acc.Balance += amount
		if err := tx.Save(&acc).Error; err != nil {
			return err
		}
		log := &entity.TransactionLog{
			ToAccountID: &acc.ID,
			Amount:      amount,
			CreatedAt:   time.Now(),
			Description: "Deposit",
		}
		return tx.Create(log).Error
	})
}

func (u *AccountUsecase) Withdraw(accountID uint, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	return u.DB.Transaction(func(tx *gorm.DB) error {
		var acc entity.Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&acc, accountID).Error; err != nil {
			return err
		}
		if acc.Balance < amount {
			return errors.New("insufficient balance")
		}
		acc.Balance -= amount
		if err := tx.Save(&acc).Error; err != nil {
			return err
		}
		log := &entity.TransactionLog{
			FromAccountID: &acc.ID,
			Amount:        amount,
			CreatedAt:     time.Now(),
			Description:   "Withdraw",
		}
		return tx.Create(log).Error
	})
}

func (u *AccountUsecase) Transfer(fromID, toID uint, amount int64) error {
	if amount <= 0 {
		return errors.New("amount must be positive")
	}
	return u.DB.Transaction(func(tx *gorm.DB) error {
		var from entity.Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&from, fromID).Error; err != nil {
			return err
		}
		var to entity.Account
		if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&to, toID).Error; err != nil {
			return err
		}
		if from.Balance < amount {
			return errors.New("insufficient balance")
		}
		from.Balance -= amount
		to.Balance += amount
		if err := tx.Save(&from).Error; err != nil {
			return err
		}
		if err := tx.Save(&to).Error; err != nil {
			return err
		}
		log := &entity.TransactionLog{
			FromAccountID: &from.ID,
			ToAccountID:   &to.ID,
			Amount:        amount,
			CreatedAt:     time.Now(),
			Description:   "Transfer",
		}
		return tx.Create(log).Error
	})
}

func (u *AccountUsecase) GetAccount(id uint) (*entity.Account, error) {
	return u.AccountRepo.FindByID(id)
}

func (u *AccountUsecase) GetTransactionLogs(accountID uint) ([]*entity.TransactionLog, error) {
	return u.TransactionRepo.FindByAccount(accountID)
}
