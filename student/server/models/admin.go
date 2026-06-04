package models

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Admin struct {
	ID        uint      `json:"id" gorm:"primary_key"`
	Username  string    `json:"username" gorm:"unique;not null"`
	Password  string    `json:"password" gorm:"not null"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
}

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
			Nickname: "管理员",
		}

		err = DB.Create(&admin).Error
		if err != nil {
			log.Fatal("Failed to create default admin: ", err)
		}
		log.Println("Default admin created")
	}
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
