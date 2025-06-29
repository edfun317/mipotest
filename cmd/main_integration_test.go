package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/edfun317/mipotest/internal/api/handler"
	"github.com/edfun317/mipotest/internal/api/router"
	"github.com/edfun317/mipotest/internal/domain/entity"
	"github.com/edfun317/mipotest/internal/infrastructure/persistence"
	"github.com/edfun317/mipotest/internal/usecase"
)

// 初始化測試 server 與 in-memory DB
func setupTestServer(t *testing.T) *gin.Engine {

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&entity.Account{}, &entity.TransactionLog{}); err != nil {
		t.Fatalf("failed to automigrate: %v", err)
	}

	accRepo := persistence.NewAccountGormRepository(db)
	txRepo := persistence.NewTransactionGormRepository(db)
	uc := &usecase.AccountUsecase{
		DB:              db,
		AccountRepo:     accRepo,
		TransactionRepo: txRepo,
	}
	h := &handler.AccountHandler{Usecase: uc}
	r := gin.Default()
	router.RegisterRoutes(r, h)
	return r
}

func TestIntegration_AccountLifecycle(t *testing.T) {
	r := setupTestServer(t)

	// 1. 建立帳戶
	createBody := map[string]interface{}{
		"name":    "Alice",
		"balance": 10000,
	}
	resp := performRequest(r, "POST", "/accounts", createBody)
	assert.Equal(t, http.StatusOK, resp.Code)
	var accResp struct {
		ID      uint   `json:"id"`
		Name    string `json:"name"`
		Balance int64  `json:"balance"`
	}
	err := json.Unmarshal(resp.Body.Bytes(), &accResp)
	assert.NoError(t, err)
	assert.Equal(t, "Alice", accResp.Name)
	assert.Equal(t, int64(10000), accResp.Balance)
	accountID := accResp.ID

	// 2. 存款
	depositBody := map[string]interface{}{
		"amount": 2500,
	}
	resp = performRequest(r, "POST", "/accounts/"+itoa(accountID)+"/deposit", depositBody)
	assert.Equal(t, http.StatusOK, resp.Code)

	// 3. 提款
	withdrawBody := map[string]interface{}{
		"amount": 1000,
	}
	resp = performRequest(r, "POST", "/accounts/"+itoa(accountID)+"/withdraw", withdrawBody)
	assert.Equal(t, http.StatusOK, resp.Code)

	// 4. 查詢帳戶
	resp = performRequest(r, "GET", "/accounts/"+itoa(accountID), nil)
	assert.Equal(t, http.StatusOK, resp.Code)
	err = json.Unmarshal(resp.Body.Bytes(), &accResp)
	assert.NoError(t, err)
	assert.Equal(t, int64(11500), accResp.Balance) // 10000+2500-1000

	// 5. 查詢交易紀錄
	resp = performRequest(r, "GET", "/accounts/"+itoa(accountID)+"/logs", nil)
	assert.Equal(t, http.StatusOK, resp.Code)
	bodyBytes, _ := io.ReadAll(resp.Body)
	var logs []map[string]interface{}
	err = json.Unmarshal(bodyBytes, &logs)
	assert.NoError(t, err)
	assert.True(t, len(logs) >= 2)
}

func TestIntegration_Transfer(t *testing.T) {
	r := setupTestServer(t)

	// 建立兩個帳戶
	var id1, id2 uint
	for _, user := range []string{"From", "To"} {
		resp := performRequest(r, "POST", "/accounts", map[string]interface{}{
			"name":    user,
			"balance": 5000,
		})
		assert.Equal(t, http.StatusOK, resp.Code)
		var respData struct {
			ID uint `json:"id"`
		}
		_ = json.Unmarshal(resp.Body.Bytes(), &respData)
		if user == "From" {
			id1 = respData.ID
		} else {
			id2 = respData.ID
		}
	}

	// 轉帳
	transferBody := map[string]interface{}{
		"from_id": id1,
		"to_id":   id2,
		"amount":  1234,
	}
	resp := performRequest(r, "POST", "/accounts/transfer", transferBody)
	assert.Equal(t, http.StatusOK, resp.Code)

	// 檢查結果
	resp = performRequest(r, "GET", "/accounts/"+itoa(id1), nil)
	var acc1 struct {
		Balance int64 `json:"balance"`
	}
	_ = json.Unmarshal(resp.Body.Bytes(), &acc1)
	assert.Equal(t, int64(3766), acc1.Balance) // 5000 - 1234

	resp = performRequest(r, "GET", "/accounts/"+itoa(id2), nil)
	var acc2 struct {
		Balance int64 `json:"balance"`
	}
	_ = json.Unmarshal(resp.Body.Bytes(), &acc2)
	assert.Equal(t, int64(6234), acc2.Balance) // 5000 + 1234
}

// 工具：模擬 http 請求
func performRequest(r http.Handler, method, path string, body interface{}) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	req := httptest.NewRequest(method, path, bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// 工具：uint 轉 string
func itoa(i uint) string {
	return strconv.FormatUint(uint64(i), 10)
}
