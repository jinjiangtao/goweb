
package handlers

import (
	"choujiang/models"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

type LotteryRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func DoLottery(c *gin.Context) {
	var req LotteryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.LotteryMutex.Lock()
	defer models.LotteryMutex.Unlock()

	var enabledPrizes []models.Prize
	models.DB.Where("enabled = ?", true).Find(&enabledPrizes)

	var availablePrizes []models.Prize
	totalProb := 0.0
	for _, p := range enabledPrizes {
		if p.StockUsed < p.Stock {
			availablePrizes = append(availablePrizes, p)
			totalProb += p.Probability
		}
	}

	var prizeName string
	isWin := false
	var selectedPrize *models.Prize = nil

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	randNum := r.Float64() * 100

	currentProb := 0.0
	for _, p := range availablePrizes {
		currentProb += p.Probability
		if randNum <= currentProb {
			prizeName = p.Name
			isWin = true
			selectedPrize = &p
			break
		}
	}

	if selectedPrize != nil {
		selectedPrize.StockUsed++
		models.DB.Save(selectedPrize)
	} else {
		prizeName = "未中奖"
	}

	status := "待领取"
	if !isWin {
		status = "已领取"
	}

	prizeID := uint(0)
	if selectedPrize != nil {
		prizeID = selectedPrize.ID
	}

	record := models.Record{
		Name:      req.Name,
		Phone:     req.Phone,
		PrizeID:   prizeID,
		PrizeName: prizeName,
		IsWin:     isWin,
		Status:    status,
		CreatedAt: time.Now(),
	}
	models.DB.Create(&record)

	c.JSON(http.StatusOK, gin.H{
		"isWin":      isWin,
		"prizeName":  prizeName,
		"record":     record,
	})
}
