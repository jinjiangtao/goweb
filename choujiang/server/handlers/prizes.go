
package handlers

import (
	"choujiang/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetPrizes(c *gin.Context) {
	var prizes []models.Prize
	models.DB.Order("id desc").Find(&amp;prizes)
	c.JSON(http.StatusOK, prizes)
}

func CreatePrize(c *gin.Context) {
	var prize models.Prize
	if err := c.ShouldBindJSON(&amp;prize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var enabledPrizes []models.Prize
	models.DB.Where("enabled = ?", true).Find(&amp;enabledPrizes)
	totalProb := prize.Probability
	for _, p := range enabledPrizes {
		totalProb += p.Probability
	}
	if totalProb &gt; 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total probability cannot exceed 100%"})
		return
	}

	prize.Enabled = true
	if err := models.DB.Create(&amp;prize).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, prize)
}

func UpdatePrize(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var prize models.Prize
	if err := models.DB.First(&amp;prize, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prize not found"})
		return
	}

	var updatedPrize models.Prize
	if err := c.ShouldBindJSON(&amp;updatedPrize); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var enabledPrizes []models.Prize
	models.DB.Where("enabled = ? AND id != ?", true, id).Find(&amp;enabledPrizes)
	totalProb := updatedPrize.Probability
	for _, p := range enabledPrizes {
		totalProb += p.Probability
	}
	if totalProb &gt; 100 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Total probability cannot exceed 100%"})
		return
	}

	prize.Name = updatedPrize.Name
	prize.Probability = updatedPrize.Probability
	prize.Stock = updatedPrize.Stock
	prize.Description = updatedPrize.Description
	prize.ImageURL = updatedPrize.ImageURL
	models.DB.Save(&amp;prize)
	c.JSON(http.StatusOK, prize)
}

func DeletePrize(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := models.DB.Delete(&amp;models.Prize{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func TogglePrize(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var prize models.Prize
	if err := models.DB.First(&amp;prize, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prize not found"})
		return
	}
	prize.Enabled = !prize.Enabled
	models.DB.Save(&amp;prize)
	c.JSON(http.StatusOK, prize)
}
