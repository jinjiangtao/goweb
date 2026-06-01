package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"goim/server/cache"
	"goim/server/handler"
	"goim/server/storage"
)

func setupRouter() *gin.Engine {
	storage.InitDB()
	cache.InitCache()
	r := gin.Default()

	r.POST("/api/auth/register", handler.Register)
	r.POST("/api/auth/login", handler.Login)

	return r
}

func TestRegister(t *testing.T) {
	r := setupRouter()

	jsonData := []byte(`{"username": "testuser", "password": "testpass", "nickname": "Test"}`)
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}
}

func TestLogin(t *testing.T) {
	r := setupRouter()

	jsonData := []byte(`{"username": "testuser", "password": "testpass"}`)
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status OK, got %d", w.Code)
	}

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	if _, ok := response["token"]; !ok {
		t.Error("Expected token in response")
	}
}