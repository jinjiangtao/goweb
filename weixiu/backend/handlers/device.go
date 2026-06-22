package handlers

import (
	"net/http"
	"strconv"
	"weixiu/config"
	"weixiu/models"

	"github.com/gin-gonic/gin"
)

type CreateDeviceRequest struct {
	DeviceCode  string            `json:"device_code" binding:"required"`
	DeviceName  string            `json:"device_name" binding:"required"`
	DeviceType  models.DeviceType `json:"device_type" binding:"required"`
	Location    string            `json:"location" binding:"required"`
	Area        string            `json:"area"`
	Status      string            `json:"status"`
	Description string            `json:"description"`
}

func CreateDevice(c *gin.Context) {
	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "detail": err.Error()})
		return
	}

	var existing models.Device
	if config.DB.Where("device_code = ?", req.DeviceCode).First(&existing).Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "设备编号已存在"})
		return
	}

	device := models.Device{
		DeviceCode:  req.DeviceCode,
		DeviceName:  req.DeviceName,
		DeviceType:  req.DeviceType,
		Location:    req.Location,
		Area:        req.Area,
		Status:      req.Status,
		Description: req.Description,
	}
	if device.Status == "" {
		device.Status = "normal"
	}

	if err := config.DB.Create(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "设备创建失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设备创建成功", "device": device})
}

func GetDevices(c *gin.Context) {
	deviceType := c.Query("device_type")
	area := c.Query("area")
	status := c.Query("status")
	keyword := c.Query("keyword")

	var devices []models.Device
	query := config.DB.Order("created_at DESC")

	if keyword != "" {
		query = query.Where("device_name LIKE ? OR device_code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	if deviceType != "" {
		query = query.Where("device_type = ?", deviceType)
	}
	if area != "" {
		query = query.Where("area = ?", area)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var total int64
	query.Model(&models.Device{}).Count(&total)
	query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&devices)

	c.JSON(http.StatusOK, gin.H{
		"list":      devices,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func GetDevice(c *gin.Context) {
	id := c.Param("id")
	var device models.Device
	if err := config.DB.First(&device, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}
	c.JSON(http.StatusOK, device)
}

func UpdateDevice(c *gin.Context) {
	id := c.Param("id")
	var device models.Device
	if err := config.DB.First(&device, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "设备不存在"})
		return
	}

	var req CreateDeviceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误", "detail": err.Error()})
		return
	}

	device.DeviceName = req.DeviceName
	device.DeviceType = req.DeviceType
	device.Location = req.Location
	device.Area = req.Area
	device.Status = req.Status
	device.Description = req.Description

	if err := config.DB.Save(&device).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "设备更新失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "设备更新成功", "device": device})
}

func DeleteDevice(c *gin.Context) {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Device{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "设备删除失败"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "设备删除成功"})
}
