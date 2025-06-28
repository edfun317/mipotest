package persistence

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/edfun317/mipotest/internal/domain/entity"
)

func setupTransactionTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory sqlite: %v", err)
	}
	// 建立資料表
	if err := db.AutoMigrate(&entity.Account{}, &entity.TransactionLog{}); err != nil {
		t.Fatalf("failed to automigrate: %v", err)
	}
	return db
}

func TestTransactionGormRepository_CreateAndFindByAccount(t *testing.T) {
	db := setupTransactionTestDB(t)
	repo := NewTransactionGormRepository(db)

	// 準備帳戶
	accA := &entity.Account{Name: "A", Balance: 1000}
	accB := &entity.Account{Name: "B", Balance: 2000}
	assert.NoError(t, db.Create(accA).Error)
	assert.NoError(t, db.Create(accB).Error)

	// 插入三筆交易（A->B、存款到A、B提款）
	log1 := &entity.TransactionLog{
		FromAccountID: &accA.ID,
		ToAccountID:   &accB.ID,
		Amount:        500,
		Description:   "Transfer",
	}
	log2 := &entity.TransactionLog{
		ToAccountID: &accA.ID,
		Amount:      300,
		Description: "Deposit",
	}
	log3 := &entity.TransactionLog{
		FromAccountID: &accB.ID,
		Amount:        100,
		Description:   "Withdraw",
	}
	assert.NoError(t, repo.Create(log1))
	assert.NoError(t, repo.Create(log2))
	assert.NoError(t, repo.Create(log3))

	// 查 accA 相關交易（應有2筆）
	logsA, err := repo.FindByAccount(accA.ID)
	assert.NoError(t, err)
	assert.Len(t, logsA, 2)
	// 檢查每一筆都至少有一個欄位是A
	for _, l := range logsA {
		fromMatch := l.FromAccountID != nil && *l.FromAccountID == accA.ID
		toMatch := l.ToAccountID != nil && *l.ToAccountID == accA.ID
		assert.True(t, fromMatch || toMatch)
	}

	// 查 accB 相關交易（應有2筆）
	logsB, err := repo.FindByAccount(accB.ID)
	assert.NoError(t, err)
	assert.Len(t, logsB, 2)
}

func TestTransactionGormRepository_FindByAccount_Empty(t *testing.T) {
	db := setupTransactionTestDB(t)
	repo := NewTransactionGormRepository(db)
	// 查詢不存在的帳戶，不應 error
	logs, err := repo.FindByAccount(9999)
	assert.NoError(t, err)
	assert.Len(t, logs, 0)
}
