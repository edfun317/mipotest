package persistence

import (
	"github.com/edfun317/mipotest/internal/domain/entity"
	"github.com/edfun317/mipotest/internal/domain/repository"

	"gorm.io/gorm"
)

type accountGormRepository struct{ db *gorm.DB }

func NewAccountGormRepository(db *gorm.DB) repository.AccountRepository {
	return &accountGormRepository{db}
}

func (r *accountGormRepository) Create(a *entity.Account) error {
	return r.db.Create(a).Error
}
func (r *accountGormRepository) FindByID(id uint) (*entity.Account, error) {
	var acc entity.Account
	if err := r.db.First(&acc, id).Error; err != nil {
		return nil, err
	}
	return &acc, nil
}
func (r *accountGormRepository) Update(a *entity.Account) error {
	return r.db.Save(a).Error
}
