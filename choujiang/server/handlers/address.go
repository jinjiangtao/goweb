package handlers

import (
	"choujiang/models"
	"encoding/json"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegionData struct {
	Name     string        `json:"name"`
	Code     string        `json:"code"`
	Children []RegionData  `json:"children"`
}

func InitAddressData(db *gorm.DB) {
	var count int64
	db.Model(&models.Address{}).Count(&count)
	if count > 0 {
		return
	}

	data, err := os.ReadFile("data/regions.json")
	if err != nil {
		return
	}

	var regions []RegionData
	if err := json.Unmarshal(data, &regions); err != nil {
		return
	}

	importRegions(db, regions, 0, 1)
}

func importRegions(db *gorm.DB, regions []RegionData, parentID uint, level int) {
	for _, region := range regions {
		addr := models.Address{
			ParentID: parentID,
			Name:     region.Name,
			Level:    level,
			Code:     region.Code,
		}
		db.Create(&addr)

		if len(region.Children) > 0 {
			importRegions(db, region.Children, addr.ID, level+1)
		}
	}
}

func GetProvinces(c *gin.Context) {
	var provinces []models.Address
	models.DB.Where("level = ?", 1).Find(&provinces)
	c.JSON(http.StatusOK, provinces)
}

func GetCities(c *gin.Context) {
	provinceID := c.Query("provinceId")
	if provinceID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "provinceId is required"})
		return
	}
	id, _ := strconv.Atoi(provinceID)
	var cities []models.Address
	models.DB.Where("parent_id = ? AND level = ?", id, 2).Find(&cities)
	c.JSON(http.StatusOK, cities)
}

func GetDistricts(c *gin.Context) {
	cityID := c.Query("cityId")
	if cityID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cityId is required"})
		return
	}
	id, _ := strconv.Atoi(cityID)
	var districts []models.Address
	models.DB.Where("parent_id = ? AND level = ?", id, 3).Find(&districts)
	c.JSON(http.StatusOK, districts)
}

func SubmitAddress(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var record models.Record
	if err := models.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	var req struct {
		ReceiverName  string `json:"receiverName"`
		ReceiverPhone string `json:"receiverPhone"`
		Province      string `json:"province"`
		City          string `json:"city"`
		District      string `json:"district"`
		DetailAddress string `json:"detailAddress"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	record.ReceiverName = req.ReceiverName
	record.ReceiverPhone = req.ReceiverPhone
	record.Province = req.Province
	record.City = req.City
	record.District = req.District
	record.DetailAddress = req.DetailAddress
	record.ShippingStatus = "未发货"
	record.Status = "已领取"

	models.DB.Save(&record)
	c.JSON(http.StatusOK, record)
}

func AdminGetProvinces(c *gin.Context) {
	var provinces []models.Address
	models.DB.Where("level = ?", 1).Find(&provinces)
	c.JSON(http.StatusOK, provinces)
}

func AdminGetCities(c *gin.Context) {
	provinceID := c.Param("provinceId")
	id, _ := strconv.Atoi(provinceID)
	var cities []models.Address
	models.DB.Where("parent_id = ? AND level = ?", id, 2).Find(&cities)
	c.JSON(http.StatusOK, cities)
}

func AdminGetDistricts(c *gin.Context) {
	cityID := c.Param("cityId")
	id, _ := strconv.Atoi(cityID)
	var districts []models.Address
	models.DB.Where("parent_id = ? AND level = ?", id, 3).Find(&districts)
	c.JSON(http.StatusOK, districts)
}

func AdminGetAddressTree(c *gin.Context) {
	var provinces []models.Address
	models.DB.Where("level = ?", 1).Find(&provinces)

	result := make([]map[string]interface{}, 0)
	for _, province := range provinces {
		var cities []models.Address
		models.DB.Where("parent_id = ? AND level = ?", province.ID, 2).Find(&cities)
		
		cityList := make([]map[string]interface{}, 0)
		for _, city := range cities {
			var districts []models.Address
			models.DB.Where("parent_id = ? AND level = ?", city.ID, 3).Find(&districts)
			
			districtList := make([]map[string]interface{}, 0)
			for _, district := range districts {
				districtList = append(districtList, map[string]interface{}{
					"id":   district.ID,
					"name": district.Name,
					"code": district.Code,
				})
			}
			
			cityList = append(cityList, map[string]interface{}{
				"id":       city.ID,
				"name":     city.Name,
				"code":     city.Code,
				"children": districtList,
			})
		}
		
		result = append(result, map[string]interface{}{
			"id":       province.ID,
			"name":     province.Name,
			"code":     province.Code,
			"children": cityList,
		})
	}

	c.JSON(http.StatusOK, result)
}

func AdminAddAddress(c *gin.Context) {
	var req struct {
		ParentID uint   `json:"parentId"`
		Name     string `json:"name"`
		Level    int    `json:"level"`
		Code     string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addr := models.Address{
		ParentID: req.ParentID,
		Name:     req.Name,
		Level:    req.Level,
		Code:     req.Code,
	}
	models.DB.Create(&addr)
	c.JSON(http.StatusOK, addr)
}

func AdminUpdateAddress(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var addr models.Address
	if err := models.DB.First(&addr, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	var req struct {
		Name string `json:"name"`
		Code string `json:"code"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	addr.Name = req.Name
	addr.Code = req.Code
	models.DB.Save(&addr)
	c.JSON(http.StatusOK, addr)
}

func AdminDeleteAddress(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var addr models.Address
	if err := models.DB.First(&addr, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Address not found"})
		return
	}

	var count int64
	models.DB.Model(&models.Address{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cannot delete address with children"})
		return
	}

	models.DB.Delete(&addr)
	c.JSON(http.StatusOK, gin.H{"message": "Deleted successfully"})
}

func AdminUpdateRecordAddress(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var record models.Record
	if err := models.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}

	var req struct {
		ReceiverName   string `json:"receiverName"`
		ReceiverPhone  string `json:"receiverPhone"`
		Province       string `json:"province"`
		City           string `json:"city"`
		District       string `json:"district"`
		DetailAddress  string `json:"detailAddress"`
		ShippingStatus string `json:"shippingStatus"`
		TrackingNumber string `json:"trackingNumber"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.ReceiverName != "" {
		record.ReceiverName = req.ReceiverName
	}
	if req.ReceiverPhone != "" {
		record.ReceiverPhone = req.ReceiverPhone
	}
	if req.Province != "" {
		record.Province = req.Province
	}
	if req.City != "" {
		record.City = req.City
	}
	if req.District != "" {
		record.District = req.District
	}
	if req.DetailAddress != "" {
		record.DetailAddress = req.DetailAddress
	}
	if req.ShippingStatus != "" {
		record.ShippingStatus = req.ShippingStatus
	}
	if req.TrackingNumber != "" {
		record.TrackingNumber = req.TrackingNumber
	}

	models.DB.Save(&record)
	c.JSON(http.StatusOK, record)
}
