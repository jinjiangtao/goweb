package test

import (
	"io"
	"net/http"
	"testing"
)

func TestOwnerAPI(t *testing.T) {
	baseURL := "http://localhost:8080/api"

	t.Run("测试获取群主列表", func(t *testing.T) {
		resp, err := http.Get(baseURL + "/owner/list")
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			t.Errorf("期望状态码 200，实际得到 %d: %s", resp.StatusCode, string(body))
		}
	})

	t.Run("测试创建群主", func(t *testing.T) {
		client := &http.Client{}
		req, err := http.NewRequest("POST", baseURL+"/admin/owner/create", nil)
		if err != nil {
			t.Fatalf("创建请求失败: %v", err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Body = io.NopCloser(io.MultiReader())

		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("请求失败: %v", err)
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		t.Logf("创建群主响应: %s", string(body))
	})
}