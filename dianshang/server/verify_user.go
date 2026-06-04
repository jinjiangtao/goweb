package main

import (
	"fmt"
	"server/models"
)

func main() {
	// 初始化数据库
	fmt.Println("初始化数据库...")
	models.InitDB()
	models.InitSuperAdmin()
	
	// 检查所有用户
	var users []models.AdminUser
	models.DB.Find(&users)
	fmt.Printf("找到 %d 个用户:\n", len(users))
	for _, u := range users {
		fmt.Printf("- ID: %d, 用户名: %s, 昵称: %s, 角色: %s, 状态: %d\n", 
			u.ID, u.Username, u.Nickname, u.Role, u.Status)
		
		// 验证密码
		if models.CheckPasswordHash("123456", u.Password) {
			fmt.Println("  ✅ 密码验证成功 (123456)")
		} else {
			fmt.Println("  ❌ 密码验证失败")
		}
	}
	
	// 模拟登录流程
	fmt.Println("\n模拟登录测试:")
	var admin models.AdminUser
	models.DB.Where("username = ?", "admin").First(&admin)
	
	if admin.Username == "" {
		fmt.Println("❌ 用户不存在")
		return
	}
	
	if admin.Status == 0 {
		fmt.Println("❌ 用户已禁用")
		return
	}
	
	if !models.CheckPasswordHash("123456", admin.Password) {
		fmt.Println("❌ 密码错误")
		return
	}
	
	fmt.Println("✅ 登录成功!")
}
