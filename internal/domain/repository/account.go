package repository

import (
	"github.com/edfun317/mipotest/internal/domain/entity"
)

type AccountRepository interface {
	Create(account *entity.Account) error
	FindByID(id uint) (*entity.Account, error)
	Update(account *entity.Account) error
}
