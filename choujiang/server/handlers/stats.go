
package handlers

import (
	"choujiang/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func GetStats(c *gin.Context) {
	var totalDraws int64
	var winCount int64
	var pendingCount int64

	models.DB.Model(&amp;models.Record{}).Count(&amp;totalDraws)
	models.DB.Model(&amp;models.Record{}).Where("is_win = ?", true).Count(&amp;winCount)
	models.DB.Model(&amp;models.Record{}).Where("status = ?", "待领取").Count(&amp;pendingCount)

	winRate := 0.0
	if totalDraws &gt; 0 {
		winRate = float64(winCount) / float64(totalDraws) * 100
	}

	type PrizeStat struct {
		PrizeName string
		Count     int64
	}
	var prizeStats []PrizeStat
	models.DB.Model(&amp;models.Record{}).Select("prize_name, count(*) as count").Where("is_win = ?", true).Group("prize_name").Scan(&amp;prizeStats)

	now := time.Now()
	var dailyStats []map[string]interface{}
	for i := 6; i &gt;= 0; i-- {
		date := now.AddDate(0, 0, -i)
		dateStr := date.Format("2006-01-02")
		var dayDraws int64
		var dayWins int64
		models.DB.Model(&amp;models.Record{}).Where("date(created_at) = ?", dateStr).Count(&amp;dayDraws)
		models.DB.Model(&amp;models.Record{}).Where("date(created_at) = ? AND is_win = ?", dateStr, true).Count(&amp;dayWins)
		dailyStats = append(dailyStats, map[string]interface{}{
			"date":  dateStr,
			"draws": dayDraws,
			"wins":  dayWins,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"totalDraws":    totalDraws,
		"winCount":      winCount,
		"winRate":       winRate,
		"pendingCount":  pendingCount,
		"prizeStats":    prizeStats,
		"dailyStats":    dailyStats,
	})
}
