package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"fangwu-server/models"
	"fangwu-server/utils"

	"github.com/gin-gonic/gin"
)

func PublicListHouses(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	keyword := c.Query("keyword")
	district := c.Query("district")
	roomType := c.Query("room_type")
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")
	sortBy := c.Query("sort_by")

	query := models.DB.Model(&models.House{}).Where("status = ?", 1)

	if keyword != "" {
		query = query.Where("title LIKE ? OR community LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if district != "" {
		query = query.Where("district = ?", district)
	}
	if roomType != "" {
		if roomType == "4室+" {
			query = query.Where("room_type LIKE ?", "%4室%")
		} else {
			query = query.Where("room_type = ?", roomType)
		}
	}
	if minPriceStr != "" {
		if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			query = query.Where("price >= ?", minPrice)
		}
	}
	if maxPriceStr != "" {
		if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			query = query.Where("price <= ?", maxPrice)
		}
	}

	switch sortBy {
	case "price_asc":
		query = query.Order("price ASC")
	case "price_desc":
		query = query.Order("price DESC")
	default:
		query = query.Order("id DESC")
	}

	var total int64
	query.Count(&total)

	var houses []models.House
	offset := (page - 1) * pageSize
	query.Offset(offset).Limit(pageSize).Find(&houses)

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
			"view_count":     h.ViewCount,
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

func PublicGetHouse(c *gin.Context) {
	id := c.Param("id")
	var house models.House
	if err := models.DB.First(&house, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "房源不存在")
		return
	}

	if house.Status != 1 {
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
		"view_count":     house.ViewCount,
		"status":         house.Status,
		"is_recommended": house.IsRecommended,
		"created_at":     house.CreatedAt,
		"updated_at":     house.UpdatedAt,
	})
}

func IncrementViewCount(c *gin.Context) {
	id := c.Param("id")
	var house models.House
	if err := models.DB.First(&house, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "房源不存在")
		return
	}

	models.DB.Model(&house).Update("view_count", house.ViewCount+1)

	utils.Success(c, gin.H{"view_count": house.ViewCount + 1})
}

func RecordConsult(c *gin.Context) {
	id := c.Param("id")
	houseID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		utils.Fail(c, http.StatusBadRequest, "无效的房源ID")
		return
	}

	var house models.House
	if err := models.DB.First(&house, houseID).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "房源不存在")
		return
	}

	consult := models.Consult{
		HouseID: uint(houseID),
	}
	models.DB.Create(&consult)

	utils.Success(c, gin.H{"success": true})
}

func GetRecommendHouses(c *gin.Context) {
	id := c.Param("id")
	var currentHouse models.House
	if err := models.DB.First(&currentHouse, id).Error; err != nil {
		utils.Fail(c, http.StatusNotFound, "房源不存在")
		return
	}

	var houses []models.House
	query := models.DB.Model(&models.House{}).Where("status = ? AND id != ?", 1, currentHouse.ID)

	if currentHouse.District != "" {
		query = query.Where("district = ?", currentHouse.District)
	}

	query.Order("is_recommended DESC, id DESC").Limit(6).Find(&houses)

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
			"room_type":      h.RoomType,
			"area":           h.Area,
			"price":          h.Price,
			"images":         images,
			"is_recommended": h.IsRecommended,
			"created_at":     h.CreatedAt,
		}
		list = append(list, item)
	}

	if list == nil {
		list = []gin.H{}
	}

	utils.Success(c, list)
}
