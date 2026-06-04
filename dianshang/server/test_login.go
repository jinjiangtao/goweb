package main

import (
	"encoding/json"
	"fmt"
	"server/handlers"
	"server/models"
	"server/utils"
)

func main() {
	// 初始化数据库
	fmt.Println("初始化数据库...")
	models.InitDB()
	models.InitSuperAdmin()
	
	// 查找admin用户
	fmt.Println("\n查找admin用户...")
	var admin models.AdminUser
	result := models.DB.Where("username = ?", "admin").First(&admin)
	if result.Error != nil {
		fmt.Printf("查询admin用户出错: %v\n", result.Error)
	} else if admin.ID == 0 {
		fmt.Println("admin用户未找到，尝试重新创建...")
		models.InitSuperAdmin()
		models.DB.Where("username = ?", "admin").First(&admin)
	}
	
	fmt.Printf("找到用户: ID=%d, 用户名=%s\n", admin.ID, admin.Username)
	
	// 测试登录
	fmt.Println("\n测试登录...")
	loginReq := handlers.LoginRequest{
		Username: "admin",
		Password: "123456",
	}
	
	// 模拟登录流程
	if admin.ID == 0 {
		fmt.Println("用户不存在，登录失败")
		return
	}
	
	if !models.CheckPasswordHash(loginReq.Password, admin.Password) {
		fmt.Println("密码错误，登录失败")
		return
	}
	
	if admin.Status == 0 {
		fmt.Println("账号被禁用")
		return
	}
	
	token, err := utils.GenerateToken(admin.ID, admin.Username, admin.Role)
	if err != nil {
		fmt.Printf("生成token失败: %v\n", err)
		return
	}
	
	fmt.Printf("登录成功！\nToken: %s\n", token)
	fmt.Println("\n完整登录响应:")
	resp := struct {
		Code int    `json:"code"`
		Msg  string `json:"message"`
		Data struct {
			Token string           `json:"token"`
			User  models.AdminUser `json:"user"`
		} `json:"data"`
	}{
		Code: 200,
		Msg:  "登录成功",
	}
	resp.Data.Token = token
	resp.Data.User = admin
	
	jsonData, _ := json.MarshalIndent(resp, "", "  ")
	fmt.Println(string(jsonData))
}
