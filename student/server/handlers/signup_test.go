package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"student-signup-server/models"
	"student-signup-server/utils"
)

func init() {
	models.InitDB()
}

func SetupTestRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/api/admin/signups", AdminCreateSignup)
	r.PUT("/api/admin/signups/:id", AdminUpdateSignup)
	r.GET("/api/admin/signups", GetSignups)
	r.GET("/api/admin/signups/export", ExportSignups)
	return r
}

func getTestToken() string {
	token, _ := utils.GenerateToken(1, "admin")
	return token
}

func TestAdminCreateSignup(t *testing.T) {
	r := SetupTestRouter()

	tests := []struct {
		name           string
		body           string
		token          string
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "创建成功",
			body:           `{"name":"测试用户","phone":"13800138000","age":18,"hukou":"北京","school":"测试学校","status":"pending"}`,
			token:          getTestToken(),
			expectedStatus: http.StatusOK,
			expectedMsg:    "创建成功",
		},
		{
			name:           "参数错误",
			body:           `{"name":"测试用户"}`,
			token:          getTestToken(),
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "参数错误",
		},
		{
			name:           "状态无效",
			body:           `{"name":"测试用户","phone":"13800138000","age":18,"hukou":"北京","school":"测试学校","status":"invalid"}`,
			token:          getTestToken(),
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "状态值无效",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/api/admin/signups", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedMsg, response["message"])
			} else {
				assert.Equal(t, tt.expectedMsg, response["error"])
			}
		})
	}
}

func TestAdminUpdateSignup(t *testing.T) {
	r := SetupTestRouter()

	signup := &models.Signup{
		Name:   "测试用户",
		Phone:  "13800138001",
		Age:    20,
		Hukou:  "上海",
		School: "测试学校",
		Status: "pending",
	}
	models.CreateSignup(signup)

	tests := []struct {
		name           string
		id             string
		body           string
		token          string
		expectedStatus int
		expectedMsg    string
	}{
		{
			name:           "更新成功",
			id:             "1",
			body:           `{"name":"更新用户","phone":"13900139000","age":21,"hukou":"广州","school":"更新学校"}`,
			token:          getTestToken(),
			expectedStatus: http.StatusOK,
			expectedMsg:    "更新成功",
		},
		{
			name:           "ID格式错误",
			id:             "invalid",
			body:           `{"name":"更新用户","phone":"13900139000","age":21,"hukou":"广州","school":"更新学校"}`,
			token:          getTestToken(),
			expectedStatus: http.StatusBadRequest,
			expectedMsg:    "ID格式错误",
		},
		{
			name:           "记录不存在",
			id:             "999",
			body:           `{"name":"更新用户","phone":"13900139000","age":21,"hukou":"广州","school":"更新学校"}`,
			token:          getTestToken(),
			expectedStatus: http.StatusNotFound,
			expectedMsg:    "记录不存在",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPut, "/api/admin/signups/"+tt.id, strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			json.Unmarshal(w.Body.Bytes(), &response)
			if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, tt.expectedMsg, response["message"])
			} else {
				assert.Equal(t, tt.expectedMsg, response["error"])
			}
		})
	}
}

func TestGetSignupsWithStatus(t *testing.T) {
	r := SetupTestRouter()

	models.CreateSignup(&models.Signup{Name: "用户A", Phone: "13800138002", Age: 18, Hukou: "北京", School: "学校A", Status: "approved"})
	models.CreateSignup(&models.Signup{Name: "用户B", Phone: "13800138003", Age: 19, Hukou: "上海", School: "学校B", Status: "pending"})

	tests := []struct {
		name           string
		query          string
		token          string
		expectedStatus int
	}{
		{
			name:           "按状态筛选approved",
			query:          "?status=approved",
			token:          getTestToken(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "按状态筛选pending",
			query:          "?status=pending",
			token:          getTestToken(),
			expectedStatus: http.StatusOK,
		},
		{
			name:           "全部状态",
			query:          "?status=all",
			token:          getTestToken(),
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/admin/signups"+tt.query, nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)
				assert.Contains(t, response, "list")
				assert.Contains(t, response, "total")
			}
		})
	}
}

func TestExportSignups(t *testing.T) {
	r := SetupTestRouter()

	req := httptest.NewRequest(http.MethodGet, "/api/admin/signups/export?status=all", nil)
	req.Header.Set("Authorization", "Bearer "+getTestToken())
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", w.Header().Get("Content-Type"))
	assert.Contains(t, w.Header().Get("Content-Disposition"), "signups.xlsx")
}