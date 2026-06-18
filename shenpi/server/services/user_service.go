package services

import (
	"errors"
	"shenpi/models"

	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) Register(user *models.User) error {
	var existingUser models.User
	result := s.db.Where("username = ?", user.Username).First(&existingUser)
	if result.Error == nil {
		return errors.New("用户名已存在")
	}
	if result.Error != gorm.ErrRecordNotFound {
		return result.Error
	}

	if err := user.HashPassword(); err != nil {
		return err
	}

	return s.db.Create(user).Error
}

func (s *UserService) Login(username, password string) (*models.User, error) {
	var user models.User
	result := s.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, errors.New("用户名或密码错误")
		}
		return nil, result.Error
	}

	if !user.CheckPassword(password) {
		return nil, errors.New("用户名或密码错误")
	}

	return &user, nil
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	var user models.User
	result := s.db.First(&user, id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (s *UserService) GetList(keyword string) ([]models.User, error) {
	var users []models.User
	query := s.db.Model(&models.User{})

	if keyword != "" {
		query = query.Where("name LIKE ? OR username LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	err := query.Order("id ASC").Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}
