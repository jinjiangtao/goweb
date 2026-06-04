package main

import (
	"fmt"
	"server/models"
)

func main() {
	fmt.Println("正在初始化数据库...")
	models.InitDB()
	fmt.Println("数据库初始化成功")
	
	fmt.Println("\n正在检查管理员用户...")
	var count int
	models.DB.Model(&models.AdminUser{}).Count(&count)
	fmt.Printf("管理员用户数量: %d\n", count)
	
	if count == 0 {
		fmt.Println("\n正在创建超级管理员...")
		models.InitSuperAdmin()
		fmt.Println("超级管理员创建成功")
	}
	
	// 列出所有管理员
	var users []models.AdminUser
	models.DB.Find(&users)
	fmt.Printf("\n找到 %d 个管理员:\n", len(users))
	for _, u := range users {
		fmt.Printf("- ID: %d, 用户名: %s, 昵称: %s, 角色: %s, 状态: %d\n", 
			u.ID, u.Username, u.Nickname, u.Role, u.Status)
	}
	
	// 尝试直接查询admin用户
	var admin models.AdminUser
	models.DB.Where("username = ?", "admin").First(&admin)
	if admin.ID > 0 {
		fmt.Printf("\n找到admin用户: ID=%d\n", admin.ID)
		
		// 验证密码
		fmt.Println("\n验证密码 (123456)...")
		if models.CheckPasswordHash("123456", admin.Password) {
			fmt.Println("密码验证成功")
		} else {
			fmt.Println("密码验证失败")
		}
	} else {
		fmt.Println("\n未找到admin用户")
	}
}
