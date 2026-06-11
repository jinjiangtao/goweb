package database

import (
	"database/sql"
	"erp/models"
	"fmt"
	"log"

	_ "modernc.org/sqlite"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	sqlDB, err := sql.Open("sqlite", "erp.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	DB, err = gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	fmt.Println("Database connected successfully")

	err = DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Menu{},
		&models.Product{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database migrated successfully")

	seedData()
}

func seedData() {
	var roleCount int64
	DB.Model(&models.Role{}).Count(&roleCount)
	if roleCount == 0 {
		adminRole := models.Role{
			Name:        "超级管理员",
			Code:        "admin",
			Description: "拥有系统所有权限",
		}
		DB.Create(&adminRole)

		menus := []models.Menu{
			{Name: "首页", Path: "/dashboard", Icon: "home", Sort: 1, Type: "menu"},
			{Name: "系统管理", Path: "/system", Icon: "setting", Sort: 2, Type: "menu"},
			{Name: "用户管理", Path: "/system/user", Icon: "user", Sort: 1, Type: "menu", ParentID: &[]uint{2}[0]},
			{Name: "角色管理", Path: "/system/role", Icon: "team", Sort: 2, Type: "menu", ParentID: &[]uint{2}[0]},
			{Name: "菜单管理", Path: "/system/menu", Icon: "menu", Sort: 3, Type: "menu", ParentID: &[]uint{2}[0]},
			{Name: "产品管理", Path: "/product", Icon: "shopping", Sort: 3, Type: "menu"},
			{Name: "产品列表", Path: "/product/list", Icon: "list", Sort: 1, Type: "menu", ParentID: &[]uint{6}[0]},
		}
		DB.Create(&menus)

		DB.Model(&adminRole).Association("Menus").Append(menus)

		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		adminUser := models.User{
			Username: "admin",
			Password: string(hashedPassword),
			Email:    "admin@example.com",
			Status:   1,
			RoleID:   adminRole.ID,
		}
		DB.Create(&adminUser)

		fmt.Println("Initial data seeded successfully")
	}
}
