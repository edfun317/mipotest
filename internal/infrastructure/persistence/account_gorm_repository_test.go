package persistence

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/edfun317/mipotest/internal/domain/entity"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}

	// 建立資料表
	if err := db.AutoMigrate(&entity.Account{}); err != nil {
		t.Fatalf("failed to automigrate: %v", err)
	}
	return db
}

func TestAccountGormRepository_CRUD(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountGormRepository(db)

	// Test Create
	acc := &entity.Account{Name: "Alice", Balance: 1000}
	err := repo.Create(acc)
	assert.NoError(t, err)
	assert.NotZero(t, acc.ID, "ID should be set after create")

	// Test FindByID
	found, err := repo.FindByID(acc.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", found.Name)
	assert.Equal(t, int64(1000), found.Balance)

	// Test Update
	found.Balance = 2000
	err = repo.Update(found)
	assert.NoError(t, err)
	updated, err := repo.FindByID(acc.ID)
	assert.NoError(t, err)
	assert.Equal(t, int64(2000), updated.Balance)
}

func TestAccountGormRepository_FindByID_NotFound(t *testing.T) {
	db := setupTestDB(t)
	repo := NewAccountGormRepository(db)

	acc, err := repo.FindByID(99999)
	assert.Nil(t, acc)
	assert.Error(t, err)
}
