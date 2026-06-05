package models

import (
	"log"

	"gorm.io/gorm"
)

type Menu struct {
	ID       uint   `json:"id" gorm:"primary_key"`
	Name     string `json:"name" gorm:"not null"`
	Path     string `json:"path"`
	Icon     string `json:"icon"`
	Sort     int    `json:"sort" gorm:"default:0"`
	ParentID uint   `json:"parent_id" gorm:"default:0"`
	Visible  int    `json:"visible" gorm:"default:1"`
	Children []Menu `json:"children" gorm:"-"`
}

func CreateDefaultMenus() {
	defaultMenus := []Menu{
		{Name: "报名列表", Path: "/", Icon: "List", Sort: 1, ParentID: 0, Visible: 1},
		{Name: "统计看板", Path: "/stats", Icon: "DataAnalysis", Sort: 2, ParentID: 0, Visible: 1},
		{Name: "用户管理", Path: "/users", Icon: "User", Sort: 3, ParentID: 0, Visible: 1},
		{Name: "菜单管理", Path: "/menus", Icon: "Menu", Sort: 4, ParentID: 0, Visible: 1},
	}

	for _, menu := range defaultMenus {
		var existing Menu
		result := DB.Where("name = ?", menu.Name).First(&existing)
		if result.Error != nil && result.Error == gorm.ErrRecordNotFound {
			err := DB.Create(&menu).Error
			if err != nil {
				log.Println("Failed to create default menu: ", err)
			}
		}
	}
}

func GetAllMenus() ([]Menu, error) {
	var menus []Menu
	err := DB.Order("sort ASC").Find(&menus).Error
	return menus, err
}

func GetMenuByRole(role string) ([]Menu, error) {
	var menus []Menu
	query := DB.Where("visible = 1").Order("sort ASC")

	if role != RoleSuperAdmin {
		query = query.Where("name NOT IN (?)", []string{"用户管理", "菜单管理"})
	}

	err := query.Find(&menus).Error
	if err != nil {
		return nil, err
	}

	return buildMenuTree(menus), nil
}

func buildMenuTree(menus []Menu) []Menu {
	menuMap := make(map[uint]*Menu)
	var rootMenus []Menu

	for i := range menus {
		menuMap[menus[i].ID] = &menus[i]
	}

	for _, menu := range menus {
		if menu.ParentID == 0 {
			rootMenus = append(rootMenus, menu)
		} else {
			if parent, ok := menuMap[menu.ParentID]; ok {
				parent.Children = append(parent.Children, menu)
			}
		}
	}

	return rootMenus
}

func CreateMenu(menu *Menu) error {
	return DB.Create(menu).Error
}

func UpdateMenu(menu *Menu) error {
	return DB.Save(menu).Error
}

func DeleteMenu(id uint) error {
	tx := DB.Begin()

	if err := tx.Delete(&Menu{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Where("parent_id = ?", id).Delete(&Menu{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func GetMenuByID(id uint) (*Menu, error) {
	var menu Menu
	err := DB.Where("id = ?", id).First(&menu).Error
	return &menu, err
}

func GetParentMenus() ([]Menu, error) {
	var menus []Menu
	err := DB.Where("parent_id = 0").Order("sort ASC").Find(&menus).Error
	return menus, err
}