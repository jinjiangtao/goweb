package database

import (
	"log"
	"richeng-server/models"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("richeng.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
	}

	err = DB.AutoMigrate(&models.Category{}, &models.Event{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	seedData()
}

func seedData() {
	var catCount int64
	DB.Model(&models.Category{}).Count(&catCount)
	var categories []models.Category
	if catCount == 0 {
		categories = []models.Category{
			{Name: "工作", Color: "#3B82F6"},
			{Name: "生活", Color: "#10B981"},
			{Name: "学习", Color: "#F59E0B"},
			{Name: "运动", Color: "#EF4444"},
			{Name: "其他", Color: "#8B5CF6"},
		}
		DB.Create(&categories)
		log.Println("Seed categories created")
	} else {
		DB.Find(&categories)
	}

	catMap := make(map[string]uint)
	for _, c := range categories {
		catMap[c.Name] = c.ID
	}

	var eventCount int64
	DB.Model(&models.Event{}).Count(&eventCount)
	if eventCount == 0 {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
		catWork := catMap["工作"]
		catLife := catMap["生活"]
		catStudy := catMap["学习"]
		catSport := catMap["运动"]

		workID := &catWork
		lifeID := &catLife
		studyID := &catStudy
		sportID := &catSport

		reminderTime := today.Add(time.Hour*8 + time.Minute*45)
		recurringEnd := today.AddDate(0, 1, 0)

		events := []models.Event{
			{
				Title:       "每日晨会",
				Description: "团队每日站会，同步工作进度",
				StartTime:   today.Add(time.Hour * 9),
				EndTime:     today.Add(time.Hour*9 + time.Minute*30),
				CategoryID:  workID,
				Location:    "会议室A",
				IsRecurring: true,
				RecurringType: "daily",
				RecurringEnd: &recurringEnd,
				HasReminder: true,
				ReminderTime: &reminderTime,
			},
			{
				Title:       "产品需求评审",
				Description: "讨论新版本功能需求和排期",
				StartTime:   today.Add(time.Hour*10 + time.Minute*30),
				EndTime:     today.Add(time.Hour*12),
				CategoryID:  workID,
				Location:    "大会议室",
			},
			{
				Title:       "午餐休息",
				Description: "午餐和午休时间",
				StartTime:   today.Add(time.Hour*12),
				EndTime:     today.Add(time.Hour*13 + time.Minute*30),
				CategoryID:  lifeID,
				Location:    "公司食堂",
			},
			{
				Title:       "Go语言学习",
				Description: "学习Gin框架和GORM最佳实践",
				StartTime:   today.Add(time.Hour*14),
				EndTime:     today.Add(time.Hour*16),
				CategoryID:  studyID,
				Location:    "工位",
			},
			{
				Title:       "健身房锻炼",
				Description: "有氧运动+力量训练",
				StartTime:   today.Add(time.Hour*18 + time.Minute*30),
				EndTime:     today.Add(time.Hour*20),
				CategoryID:  sportID,
				Location:    "商业健身房",
			},
			{
				Title:       "代码评审",
				Description: "评审团队成员的Pull Request",
				StartTime:   today.AddDate(0, 0, 1).Add(time.Hour * 10),
				EndTime:     today.AddDate(0, 0, 1).Add(time.Hour*11 + time.Minute*30),
				CategoryID:  workID,
				Location:    "线上会议",
			},
			{
				Title:       "每周周会",
				Description: "团队周度总结与计划会议",
				StartTime:   today.AddDate(0, 0, 2).Add(time.Hour * 15),
				EndTime:     today.AddDate(0, 0, 2).Add(time.Hour*16 + time.Minute*30),
				CategoryID:  workID,
				Location:    "会议室B",
				IsRecurring: true,
				RecurringType: "weekly",
			},
		}
		DB.Create(&events)
		log.Println("Seed events created:", len(events))
	}
}
