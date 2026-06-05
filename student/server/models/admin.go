package models

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID            uint      `json:"id" gorm:"primary_key"`
	Username      string    `json:"username" gorm:"unique;not null"`
	Password      string    `json:"password" gorm:"not null"`
	Nickname      string    `json:"nickname"`
	Role          string    `json:"role" gorm:"default:admin"`
	Status        int       `json:"status" gorm:"default:1"`
	LastLoginTime time.Time `json:"last_login_time"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

const (
	RoleSuperAdmin = "super_admin"
	RoleAdmin      = "admin"
	StatusEnabled  = 1
	StatusDisabled = 0
)

func CreateDefaultAdmin() {
	var admin Admin
	result := DB.Where("username = ?", "admin").First(&admin)
	if result.Error != nil && result.Error.Error() == "record not found" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
		if err != nil {
			log.Fatal("Failed to hash password: ", err)
		}

		admin := Admin{
			Username: "admin",
			Password: string(hashedPassword),
			Nickname: "超级管理员",
			Role:     RoleSuperAdmin,
			Status:   StatusEnabled,
		}

		err = DB.Create(&admin).Error
		if err != nil {
			log.Fatal("Failed to create default admin: ", err)
		}
		log.Println("Default admin created")
	}
}

func GetAdminByID(id uint) (*Admin, error) {
	var admin Admin
	err := DB.Where("id = ?", id).First(&admin).Error
	return &admin, err
}

func GetAdmins(username string, page, pageSize int) ([]Admin, int64, error) {
	var admins []Admin
	var total int64

	query := DB.Model(&Admin{})

	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Order("created_at DESC").Offset((page - 1) * pageSize).Limit(pageSize).Find(&admins).Error
	return admins, total, err
}

func CreateAdmin(admin *Admin) error {
	return DB.Create(admin).Error
}

func UpdateAdmin(admin *Admin) error {
	return DB.Save(admin).Error
}

func DeleteAdmin(id uint) error {
	return DB.Delete(&Admin{}, id).Error
}

func GetAdminByUsername(username string) (*Admin, error) {
	var admin Admin
	err := DB.Where("username = ?", username).First(&admin).Error
	return &admin, err
}

func (admin *Admin) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	return err == nil
}
