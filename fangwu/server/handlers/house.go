package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"fangwu-server/models"
	"fangwu-server/utils"

	"github.com/gin-gonic/gin"
)

type HouseRequest struct {
	Title        string  `json:"title" binding:"required"`
	Community    string  `json:"community" binding:"required"`
	District     string  `json:"district" binding:"required"`
	Address      string  `json:"address" binding:"required"`
	RoomType     string  `json:"room_type" binding:"required"`
	Area         float64 `json:"area" binding:"required"`
	Price        float64 `json:"price" binding:"required"`
	Floor        string  `json:"floor"`
	Orientation  string  `json:"orientation"`
	Decoration   string  `json:"decoration"`
	Facilities   string  `json:"facilities"`
	Description  string  `json:"description" binding:"required"`
	ContactName  string  `json:"contact_name" binding:"required"`
	ContactPhone string  `json:"contact_phone" binding:"required"`
	Images       string  `json:"images"`
}

func ListHouses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	title := c.Query("title")
	community := c.Query("community")
	status := c.Query("status")
	isRecommended := c.Query("is_recommended")

	query := models.DB.Model(&models.House{})

	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	if community != "" {
		query = query.Where("community LIKE ?", "%"+community+"%")
	}
	if status != "" {
		s, _ := strconv.Atoi(status)
		query = query.Where("status = ?", s)
	}
	if isRecommended != "" {
		r := isRecommended == "1" || strings.ToLower(isRecommended) == "true"
		query = query.Where("is_recommended = ?", r)
	}

	var total int64
	query.Count(&total)

	var houses []models.House
	offset := (page - 1) * pageSize
	query.Order("id DESC").Offset(offset).Limit(pageSize).Find(&houses)

	var list []gin.H
	for _, h := range houses {
		var images []string
		if h.Images != "" {
			json.Unmarshal([]byte(h.Images), &images)
		}
		item := gin.H{
			"id":             h.ID,
			"title":          h.Title,
			"community":      h.Community,
			"district":       h.District,
			"address":        h.Address,
			"room_type":      h.RoomType,
			"area":           h.Area,
			"price":          h.Price,
			"floor":          h.Floor,
			"orientation":    h.Orientation,
			"decoration":     h.Decoration,
			"facilities":     h.Facilities,
			"description":    h.Description,
			"contact_name":   h.ContactName,
			"contact_phone":  h.ContactPhone,
			"images":         images,
			"status":         h.Status,
			"is_recommended": h.IsRecommended,
			"created_at":     h.CreatedAt,
			"updated_at":     h.UpdatedAt,
		}
		list = append(list, item)
	}

	if list == nil {
		list = []gin.H{}
	}

	utils.PageResult(c, list, total, page, pageSize)
}

func GetHouse(c *gin.Context) {
	id := c.Param("id")
	var house models.House
	if err := models.DB.First(&house, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "房源不存在")
		return
	}

	var images []string
	if house.Images != "" {
		json.Unmarshal([]byte(house.Images), &images)
	}

	utils.Success(c, gin.H{
		"id":             house.ID,
		"title":          house.Title,
		"community":      house.Community,
		"district":       house.District,
		"address":        house.Address,
		"room_type":      house.RoomType,
		"area":           house.Area,
		"price":          house.Price,
		"floor":          house.Floor,
		"orientation":    house.Orientation,
		"decoration":     house.Decoration,
		"facilities":     house.Facilities,
		"description":    house.Description,
		"contact_name":   house.ContactName,
		"contact_phone":  house.ContactPhone,
		"images":         images,
		"status":         house.Status,
		"is_recommended": house.IsRecommended,
		"created_at":     house.CreatedAt,
		"updated_at":     house.UpdatedAt,
	})
}

func CreateHouse(c *gin.Context) {
	var req HouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "请填写必填字段")
		return
	}

	house := models.House{
		Title:        req.Title,
		Community:    req.Community,
		District:     req.District,
		Address:      req.Address,
		RoomType:     req.RoomType,
		Area:         req.Area,
		Price:        req.Price,
		Floor:        req.Floor,
		Orientation:  req.Orientation,
		Decoration:   req.Decoration,
		Facilities:   req.Facilities,
		Description:  req.Description,
		ContactName:  req.ContactName,
		ContactPhone: req.ContactPhone,
		Images:       req.Images,
		Status:       1,
	}

	if err := models.DB.Create(&house).Error; err != nil {
		utils.Fail(c, http.StatusInternalServerError, "创建房源失败")
		return
	}

	utils.Success(c, gin.H{"id": house.ID})
}

func UpdateHouse(c *gin.Context) {
	id := c.Param("id")
	var house models.House
	if err := models.DB.First(&house, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "房源不存在")
		return
	}

	var req HouseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Fail(c, http.StatusBadRequest, "请填写必填字段")
		return
	}

	models.DB.Model(&house).Updates(map[string]interface{}{
		"title":         req.Title,
		"community":     req.Community,
		"district":      req.District,
		"address":       req.Address,
		"room_type":     req.RoomType,
		"area":          req.Area,
		"price":         req.Price,
		"floor":         req.Floor,
		"orientation":   req.Orientation,
		"decoration":    req.Decoration,
		"facilities":    req.Facilities,
		"description":   req.Description,
		"contact_name":  req.ContactName,
		"contact_phone": req.ContactPhone,
		"images":        req.Images,
	})

	utils.Success(c, nil)
}

func DeleteHouse(c *gin.Context) {
	id := c.Param("id")
	var house models.House
	if err := models.DB.First(&house, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "房源不存在")
		return
	}
	models.DB.Delete(&house)
	utils.Success(c, nil)
}

func ToggleStatus(c *gin.Context) {
	id := c.Param("id")
	var house models.House
	if err := models.DB.First(&house, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "房源不存在")
		return
	}

	newStatus := 1
	if house.Status == 1 {
		newStatus = 0
	}

	models.DB.Model(&house).Update("status", newStatus)
	utils.Success(c, gin.H{"status": newStatus})
}

func ToggleRecommend(c *gin.Context) {
	id := c.Param("id")
	var house models.House
	if err := models.DB.First(&house, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "房源不存在")
		return
	}

	models.DB.Model(&house).Update("is_recommended", !house.IsRecommended)
	utils.Success(c, gin.H{"is_recommended": !house.IsRecommended})
}
