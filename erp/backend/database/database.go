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
		&models.Customer{},
		&models.Supplier{},
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

		// 创建菜单
		menu1 := models.Menu{Name: "首页", Path: "/dashboard", Icon: "House", Sort: 1, Hidden: false}
		menu2 := models.Menu{Name: "系统设置", Path: "", Icon: "Setting", Sort: 2, Hidden: false}
		menu3 := models.Menu{Name: "用户管理", Path: "/system/user", Icon: "User", Sort: 1, Hidden: false}
		menu4 := models.Menu{Name: "角色管理", Path: "/system/role", Icon: "UserFilled", Sort: 2, Hidden: false}
		menu5 := models.Menu{Name: "菜单管理", Path: "/system/menu", Icon: "Menu", Sort: 3, Hidden: false}
		menu6 := models.Menu{Name: "产品管理", Path: "/product", Icon: "Goods", Sort: 3, Hidden: false}
		menu7 := models.Menu{Name: "客户管理", Path: "/customer", Icon: "User", Sort: 4, Hidden: false}
		menu8 := models.Menu{Name: "供应商管理", Path: "/supplier", Icon: "OfficeBuilding", Sort: 5, Hidden: false}

		DB.Create(&menu1)
		DB.Create(&menu2)
		DB.Create(&menu3)
		DB.Create(&menu4)
		DB.Create(&menu5)
		DB.Create(&menu6)
		DB.Create(&menu7)
		DB.Create(&menu8)

		// 设置父子关系
		menu3.ParentID = &menu2.ID
		menu4.ParentID = &menu2.ID
		menu5.ParentID = &menu2.ID
		DB.Save(&menu3)
		DB.Save(&menu4)
		DB.Save(&menu5)

		menus := []models.Menu{menu1, menu2, menu3, menu4, menu5, menu6, menu7, menu8}

		// 为角色分配菜单
		for _, menu := range menus {
			DB.Exec("INSERT INTO role_menus (role_id, menu_id) VALUES (?, ?)", adminRole.ID, menu.ID)
		}

		// 创建默认产品
		products := []models.Product{
			{Name: "产品1", Code: "PROD001", Price: 199.99, Spec: "产品1规格"},
			{Name: "产品2", Code: "PROD002", Price: 299.99, Spec: "产品2规格"},
			{Name: "产品3", Code: "PROD003", Price: 399.99, Spec: "产品3规格"},
		}
		for _, p := range products {
			DB.Create(&p)
		}

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
