package db

import (
	"jianli-server/models"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.New(sqlite.Config{
		DSN: "file:jianli.db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)",
		DriverName: "sqlite",
	}), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = DB.AutoMigrate(&models.Template{}, &models.Resume{}, &models.Version{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	seedTemplates()

	log.Println("Database initialized successfully")
}

func seedTemplates() {
	var count int64
	DB.Model(&models.Template{}).Count(&count)
	if count > 0 {
		return
	}

	templates := []models.Template{
		{
			Name:        "简约商务模板",
			Description: "简洁大方的商务风格，适合求职各类岗位",
			Category:    "商务",
			Thumbnail:   "",
			Content: models.JSON{
				"modules": []map[string]interface{}{
					{"id": "basic", "name": "基本信息", "visible": true, "order": 1},
					{"id": "education", "name": "教育背景", "visible": true, "order": 2},
					{"id": "experience", "name": "工作经历", "visible": true, "order": 3},
					{"id": "skills", "name": "专业技能", "visible": true, "order": 4},
					{"id": "projects", "name": "项目经验", "visible": true, "order": 5},
				},
				"basic": map[string]interface{}{
					"name": "姓名", "phone": "手机号", "email": "邮箱", "address": "所在城市",
				},
				"education": []map[string]interface{}{},
				"experience": []map[string]interface{}{},
				"skills": []map[string]interface{}{},
				"projects": []map[string]interface{}{},
			},
			StyleConfig: models.JSON{
				"font_family":    "Microsoft YaHei",
				"font_size":      12,
				"primary_color":  "#1a1a2e",
				"secondary_color": "#16213e",
				"accent_color":   "#0f3460",
				"line_spacing":   1.5,
				"margin_top":     20,
				"margin_bottom":  20,
				"margin_left":    25,
				"margin_right":   25,
			},
		},
		{
			Name:        "创意设计模板",
			Description: "充满创意的设计风格，适合设计、创意类岗位",
			Category:    "创意",
			Thumbnail:   "",
			Content: models.JSON{
				"modules": []map[string]interface{}{
					{"id": "basic", "name": "基本信息", "visible": true, "order": 1},
					{"id": "summary", "name": "个人简介", "visible": true, "order": 2},
					{"id": "skills", "name": "专业技能", "visible": true, "order": 3},
					{"id": "experience", "name": "工作经历", "visible": true, "order": 4},
					{"id": "projects", "name": "作品集", "visible": true, "order": 5},
					{"id": "education", "name": "教育背景", "visible": true, "order": 6},
				},
				"basic": map[string]interface{}{
					"name": "姓名", "title": "职位", "phone": "手机号", "email": "邮箱", "website": "个人网站",
				},
				"summary":    "",
				"education":  []map[string]interface{}{},
				"experience": []map[string]interface{}{},
				"skills":     []map[string]interface{}{},
				"projects":   []map[string]interface{}{},
			},
			StyleConfig: models.JSON{
				"font_family":    "Microsoft YaHei",
				"font_size":      12,
				"primary_color":  "#2c3e50",
				"secondary_color": "#34495e",
				"accent_color":   "#e74c3c",
				"line_spacing":   1.6,
				"margin_top":     15,
				"margin_bottom":  15,
				"margin_left":    20,
				"margin_right":   20,
			},
		},
		{
			Name:        "技术工程师模板",
			Description: "专业技术风格，突出技能和项目经验",
			Category:    "技术",
			Thumbnail:   "",
			Content: models.JSON{
				"modules": []map[string]interface{}{
					{"id": "basic", "name": "基本信息", "visible": true, "order": 1},
					{"id": "summary", "name": "技术简介", "visible": true, "order": 2},
					{"id": "skills", "name": "技术栈", "visible": true, "order": 3},
					{"id": "projects", "name": "项目经验", "visible": true, "order": 4},
					{"id": "experience", "name": "工作经历", "visible": true, "order": 5},
					{"id": "education", "name": "教育背景", "visible": true, "order": 6},
				},
				"basic": map[string]interface{}{
					"name": "姓名", "phone": "手机号", "email": "邮箱", "github": "GitHub", "location": "所在城市",
				},
				"summary":    "",
				"education":  []map[string]interface{}{},
				"experience": []map[string]interface{}{},
				"skills":     []map[string]interface{}{},
				"projects":   []map[string]interface{}{},
			},
			StyleConfig: models.JSON{
				"font_family":    "Consolas",
				"font_size":      11,
				"primary_color":  "#000000",
				"secondary_color": "#333333",
				"accent_color":   "#0066cc",
				"line_spacing":   1.4,
				"margin_top":     20,
				"margin_bottom":  20,
				"margin_left":    25,
				"margin_right":   25,
			},
		},
	}

	for _, t := range templates {
		DB.Create(&t)
	}

	log.Println("Template data seeded")
}
