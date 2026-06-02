package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"goim/server/handler"
	"goim/server/storage"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupOwnerRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/api/admin/owner/create", handler.CreateOwner)
	r.GET("/api/owner/list", handler.GetAllOwners)
	r.GET("/api/owner/:ownerID", handler.GetOwnerByID)
	r.POST("/api/owner/join", handler.JoinOwner)
	r.POST("/api/owner/leave", handler.LeaveOwner)
	r.GET("/api/owner/members/:ownerID", handler.GetOwnerMembers)
	r.DELETE("/api/admin/owner/:ownerID", handler.DeleteOwner)
	r.POST("/api/admin/owner/remove-member", handler.RemoveOwnerMember)
	return r
}

func TestCreateOwner(t *testing.T) {
	storage.InitDB()
	r := setupOwnerRouter()

	jsonData := []byte(`{"name":"测试群主","description":"这是一个测试群主","avatar":""}`)
	req, _ := http.NewRequest("POST", "/api/admin/owner/create", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.NotEmpty(t, response["id"])
	assert.Equal(t, "测试群主", response["name"])
	assert.Equal(t, "这是一个测试群主", response["description"])
}

func TestGetAllOwners(t *testing.T) {
	storage.InitDB()
	r := setupOwnerRouter()

	req, _ := http.NewRequest("GET", "/api/owner/list", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response []interface{}
	json.Unmarshal(w.Body.Bytes(), &response)
	assert.IsType(t, []interface{}{}, response)
}

func TestGetOwnerByID(t *testing.T) {
	storage.InitDB()
	r := setupOwnerRouter()

	req, _ := http.NewRequest("POST", "/api/admin/owner/create", bytes.NewBuffer([]byte(`{"name":"测试群主2","description":"描述","avatar":""}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	ownerID := createResp["id"].(string)

	req2, _ := http.NewRequest("GET", "/api/owner/"+ownerID, nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	var response map[string]interface{}
	json.Unmarshal(w2.Body.Bytes(), &response)
	assert.Equal(t, "测试群主2", response["name"])
}

func TestJoinAndLeaveOwner(t *testing.T) {
	storage.InitDB()
	r := setupOwnerRouter()

	req, _ := http.NewRequest("POST", "/api/admin/owner/create", bytes.NewBuffer([]byte(`{"name":"测试群主3","description":"描述","avatar":""}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	ownerID := createResp["id"].(string)

	joinData := []byte(`{"owner_id":"` + ownerID + `","user_id":"test_user_123"}`)
	req2, _ := http.NewRequest("POST", "/api/owner/join", bytes.NewBuffer(joinData))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	req3, _ := http.NewRequest("GET", "/api/owner/members/"+ownerID, nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusOK, w3.Code)
	var members []interface{}
	json.Unmarshal(w3.Body.Bytes(), &members)
	assert.Len(t, members, 1)

	leaveData := []byte(`{"owner_id":"` + ownerID + `","user_id":"test_user_123"}`)
	req4, _ := http.NewRequest("POST", "/api/owner/leave", bytes.NewBuffer(leaveData))
	req4.Header.Set("Content-Type", "application/json")
	w4 := httptest.NewRecorder()
	r.ServeHTTP(w4, req4)

	assert.Equal(t, http.StatusOK, w4.Code)

	req5, _ := http.NewRequest("GET", "/api/owner/members/"+ownerID, nil)
	w5 := httptest.NewRecorder()
	r.ServeHTTP(w5, req5)

	assert.Equal(t, http.StatusOK, w5.Code)
	var membersAfter []interface{}
	json.Unmarshal(w5.Body.Bytes(), &membersAfter)
	assert.Len(t, membersAfter, 0)
}

func TestDeleteOwner(t *testing.T) {
	storage.InitDB()
	r := setupOwnerRouter()

	req, _ := http.NewRequest("POST", "/api/admin/owner/create", bytes.NewBuffer([]byte(`{"name":"测试群主4","description":"描述","avatar":""}`)))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	var createResp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &createResp)
	ownerID := createResp["id"].(string)

	req2, _ := http.NewRequest("DELETE", "/api/admin/owner/"+ownerID, nil)
	w2 := httptest.NewRecorder()
	r.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	req3, _ := http.NewRequest("GET", "/api/owner/"+ownerID, nil)
	w3 := httptest.NewRecorder()
	r.ServeHTTP(w3, req3)

	assert.Equal(t, http.StatusNotFound, w3.Code)
}