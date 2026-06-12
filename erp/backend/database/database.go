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
		&models.PurchaseOrder{},
		&models.PurchaseOrderItem{},
		&models.SalesOrder{},
		&models.SalesOrderItem{},
	)
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	fmt.Println("Database migrated successfully")

	seedData()
}

func seedData() {
	// 获取或创建超级管理员角色
	var adminRole models.Role
	DB.Where("code = ?", "admin").FirstOrCreate(&adminRole, models.Role{
		Name:        "超级管理员",
		Code:        "admin",
		Description: "拥有系统所有权限",
	})

	// 定义所有需要的菜单
	menuDefinitions := []models.Menu{
		{Name: "首页", Path: "/dashboard", Icon: "House", Sort: 1, Hidden: false},
		{Name: "系统设置", Path: "", Icon: "Setting", Sort: 2, Hidden: false},
		{Name: "用户管理", Path: "/system/user", Icon: "User", Sort: 1, Hidden: false},
		{Name: "角色管理", Path: "/system/role", Icon: "UserFilled", Sort: 2, Hidden: false},
		{Name: "菜单管理", Path: "/system/menu", Icon: "Menu", Sort: 3, Hidden: false},
		{Name: "产品管理", Path: "/product", Icon: "Goods", Sort: 3, Hidden: false},
		{Name: "客户管理", Path: "/customer", Icon: "User", Sort: 4, Hidden: false},
		{Name: "供应商管理", Path: "/supplier", Icon: "OfficeBuilding", Sort: 5, Hidden: false},
		{Name: "采购订单", Path: "/purchase-order", Icon: "ShoppingCart", Sort: 6, Hidden: false},
		{Name: "销售订单", Path: "/sales-order", Icon: "Sell", Sort: 7, Hidden: false},
	}

	var createdMenus []models.Menu
	for _, menuDef := range menuDefinitions {
		var existingMenu models.Menu
		// 使用Path作为唯一标识来查找或创建菜单
		result := DB.Where("path = ?", menuDef.Path).FirstOrCreate(&existingMenu, menuDef)
		if result.Error == nil {
			createdMenus = append(createdMenus, existingMenu)
		}
	}

	// 设置父子关系：系统设置是父菜单
	var systemMenu, userMenu, roleMenu, menuMenu models.Menu
	DB.Where("path = ?", "").First(&systemMenu)
	DB.Where("path = ?", "/system/user").First(&userMenu)
	DB.Where("path = ?", "/system/role").First(&roleMenu)
	DB.Where("path = ?", "/system/menu").First(&menuMenu)

	if systemMenu.ID > 0 {
		userMenu.ParentID = &systemMenu.ID
		roleMenu.ParentID = &systemMenu.ID
		menuMenu.ParentID = &systemMenu.ID
		DB.Save(&userMenu)
		DB.Save(&roleMenu)
		DB.Save(&menuMenu)
	}

	// 为角色分配菜单（先清空再添加，确保所有菜单都有）
	DB.Exec("DELETE FROM role_menus WHERE role_id = ?", adminRole.ID)
	for _, menu := range createdMenus {
		DB.Exec("INSERT OR IGNORE INTO role_menus (role_id, menu_id) VALUES (?, ?)", adminRole.ID, menu.ID)
	}

	// 创建默认产品（如果不存在）
	var productCount int64
	DB.Model(&models.Product{}).Count(&productCount)
	if productCount == 0 {
		products := []models.Product{
			{Name: "产品1", Code: "PROD001", Price: 199.99, Spec: "产品1规格"},
			{Name: "产品2", Code: "PROD002", Price: 299.99, Spec: "产品2规格"},
			{Name: "产品3", Code: "PROD003", Price: 399.99, Spec: "产品3规格"},
		}
		for _, p := range products {
			DB.FirstOrCreate(&p, models.Product{Code: p.Code})
		}
	}

	// 创建默认管理员用户（如果不存在）
	var userCount int64
	DB.Model(&models.User{}).Where("username = ?", "admin").Count(&userCount)
	if userCount == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		adminUser := models.User{
			Username: "admin",
			Password: string(hashedPassword),
			Email:    "admin@example.com",
			Status:   1,
			RoleID:   adminRole.ID,
		}
		DB.Create(&adminUser)
	}

	fmt.Println("Initial data seeded/updated successfully")
}
