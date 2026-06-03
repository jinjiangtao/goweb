package models

import (
	"crypto/md5"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
	_ "modernc.org/sqlite"
	"golang.org/x/crypto/bcrypt"
)

var DB *gorm.DB

type AdminUser struct {
	ID           uint      `gorm:"primary_key" json:"id"`
	Username     string    `gorm:"unique;not null" json:"username"`
	Password     string    `gorm:"not null" json:"-"`
	Nickname     string    `json:"nickname"`
	Role         string    `gorm:"not null" json:"role"`
	Status       int       `gorm:"default:1" json:"status"`
	LastLoginAt  time.Time `json:"last_login_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type Menu struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	ParentID  uint      `json:"parent_id"`
	Name      string    `gorm:"not null" json:"name"`
	Path      string    `json:"path"`
	Icon      string    `json:"icon"`
	Sort      int       `gorm:"default:0" json:"sort"`
	Visible   int       `gorm:"default:1" json:"visible"`
	CreatedAt time.Time `json:"created_at"`
	Children  []Menu    `gorm:"-" json:"children"`
}

type RoleMenu struct {
	ID     uint   `gorm:"primary_key" json:"id"`
	Role   string `gorm:"not null" json:"role"`
	MenuID uint   `gorm:"not null" json:"menu_id"`
}

func InitDB() {
	var err error
	DB, err = gorm.Open("sqlite", "file:ecommerce.db?mode=rwc")
	if err != nil {
		panic("failed to connect database: " + err.Error())
	}
	DB.AutoMigrate(&AdminUser{}, &Menu{}, &RoleMenu{})
}

func InitSuperAdmin() {
	var count int
	DB.Model(&AdminUser{}).Count(&count)
	if count == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		superAdmin := AdminUser{
			Username: "admin",
			Password: string(hashedPassword),
			Nickname: "超级管理员",
			Role:     "super",
			Status:   1,
		}
		DB.Create(&superAdmin)
		InitDefaultMenus()
	}
}

func InitDefaultMenus() {
	menus := []Menu{
		{Name: "管理员管理", Path: "/admin/users", Icon: "user", Sort: 1, Visible: 1, ParentID: 0},
		{Name: "菜单管理", Path: "/admin/menus", Icon: "menu", Sort: 2, Visible: 1, ParentID: 0},
		{Name: "角色权限", Path: "/admin/roles", Icon: "lock", Sort: 3, Visible: 1, ParentID: 0},
	}
	for _, menu := range menus {
		DB.Create(&menu)
	}
	AssignAllMenusToSuper()
}

func AssignAllMenusToSuper() {
	var menus []Menu
	DB.Find(&menus)
	for _, menu := range menus {
		roleMenu := RoleMenu{Role: "super", MenuID: menu.ID}
		DB.Create(&roleMenu)
	}
}

func GeneratePasswordHash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateMD5(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}