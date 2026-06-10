
package handlers

import (
	"choujiang/models"
	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
	"net/http"
	"strconv"
)

func GetRecords(c *gin.Context) {
	name := c.Query("name")
	phone := c.Query("phone")
	prizeName := c.Query("prizeName")
	isWin := c.Query("isWin")

	query := models.DB.Order("id desc")
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	if prizeName != "" {
		query = query.Where("prize_name = ?", prizeName)
	}
	if isWin != "" {
		if isWin == "true" {
			query = query.Where("is_win = ?", true)
		} else {
			query = query.Where("is_win = ?", false)
		}
	}

	var records []models.Record
	query.Find(&records)
	c.JSON(http.StatusOK, records)
}

func GetRecordsByPhone(c *gin.Context) {
	phone := c.Query("phone")
	if phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone is required"})
		return
	}

	var records []models.Record
	models.DB.Where("phone = ?", phone).Order("id desc").Find(&records)
	c.JSON(http.StatusOK, records)
}

func ClaimRecord(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var record models.Record
	if err := models.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	record.Status = "已领取"
	models.DB.Save(&record)
	c.JSON(http.StatusOK, record)
}

func ClaimRecordPublic(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var record models.Record
	if err := models.DB.First(&record, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Record not found"})
		return
	}
	record.Status = "已领取"
	models.DB.Save(&record)
	c.JSON(http.StatusOK, record)
}

func ExportRecords(c *gin.Context) {
	var records []models.Record
	models.DB.Order("id desc").Find(&records)

	f := excelize.NewFile()
	sheetName := "抽奖记录"
	f.SetSheetName("Sheet1", sheetName)

	headers := []string{"ID", "姓名", "手机号", "奖品", "是否中奖", "状态", "时间"}
	for i, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	for i, record := range records {
		row := i + 2
		f.SetCellValue(sheetName, "A"+strconv.Itoa(row), record.ID)
		f.SetCellValue(sheetName, "B"+strconv.Itoa(row), record.Name)
		f.SetCellValue(sheetName, "C"+strconv.Itoa(row), record.Phone)
		f.SetCellValue(sheetName, "D"+strconv.Itoa(row), record.PrizeName)
		isWinStr := "否"
		if record.IsWin {
			isWinStr = "是"
		}
		f.SetCellValue(sheetName, "E"+strconv.Itoa(row), isWinStr)
		f.SetCellValue(sheetName, "F"+strconv.Itoa(row), record.Status)
		f.SetCellValue(sheetName, "G"+strconv.Itoa(row), record.CreatedAt.Format("2006-01-02 15:04:05"))
	}

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", "attachment; filename=records.xlsx")
	f.Write(c.Writer)
}

