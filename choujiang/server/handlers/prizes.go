
package handlers

import (
	"choujiang/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetPrizes(c *gin.Context) {
	var prizes []models.Prize
	models.DB.Order("id desc").Find(&prizes)
	c.JSON(http.StatusOK, prizes)
}

func CreatePrize(c *gin.Context) {
	var prize models.Prize
	if err := c.ShouldBindJSON(&prize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var enabledPrizes []models.Prize
	models.DB.Where("enabled = ?", true).Find(&enabledPrizes)
	totalProb := prize.Probability
	for _, p := range enabledPrizes {
		totalProb += p.Probability
	}
	if totalProb > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total probability cannot exceed 100%"})
		return
	}

	prize.Enabled = true
	if err := models.DB.Create(&prize).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prize)
}

func UpdatePrize(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var prize models.Prize
	if err := models.DB.First(&prize, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prize not found"})
		return
	}

	var updatedPrize models.Prize
	if err := c.ShouldBindJSON(&updatedPrize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var enabledPrizes []models.Prize
	models.DB.Where("enabled = ? AND id != ?", true, id).Find(&enabledPrizes)
	totalProb := updatedPrize.Probability
	for _, p := range enabledPrizes {
		totalProb += p.Probability
	}
	if totalProb > 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total probability cannot exceed 100%"})
		return
	}

	prize.Name = updatedPrize.Name
	prize.Probability = updatedPrize.Probability
	prize.Stock = updatedPrize.Stock
	prize.Description = updatedPrize.Description
	prize.ImageURL = updatedPrize.ImageURL
	models.DB.Save(&prize)
	c.JSON(http.StatusOK, prize)
}

func DeletePrize(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := models.DB.Delete(&models.Prize{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func TogglePrize(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var prize models.Prize
	if err := models.DB.First(&prize, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prize not found"})
		return
	}
	prize.Enabled = !prize.Enabled
	models.DB.Save(&prize)
	c.JSON(http.StatusOK, prize)
}
