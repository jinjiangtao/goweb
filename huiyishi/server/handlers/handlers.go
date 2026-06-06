
package handlers

import (
	"net/http"
	"strconv"
	"time"
	"huiyishi-server/database"
	"huiyishi-server/models"
	"huiyishi-server/utils"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var admin models.Admin
	if err := database.DB.Where("username = ?", req.Username).First(&admin).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	token, err := utils.GenerateToken(admin.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"admin": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"nickname": admin.Nickname,
		},
	})
}

func GetRooms(c *gin.Context) {
	search := c.Query("search")

	var rooms []models.Room
	query := database.DB.Order("id desc")

	if search != "" {
		query = query.Where("name LIKE ?", "%"+search+"%")
	}

	query.Find(&rooms)
	c.JSON(http.StatusOK, rooms)
}

func CreateRoom(c *gin.Context) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}
	room.CreatedAt = time.Now()
	database.DB.Create(&room)
	c.JSON(http.StatusOK, room)
}

func UpdateRoom(c *gin.Context) {
	id := c.Param("id")
	var room models.Room
	if err := database.DB.First(&room, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "会议室不存在"})
		return
	}

	var updateData models.Room
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	database.DB.Model(&room).Updates(updateData)
	c.JSON(http.StatusOK, room)
}

func DeleteRoom(c *gin.Context) {
	id := c.Param("id")
	var count int64
	database.DB.Model(&models.Booking{}).Where("room_id = ?", id).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该会议室有预订记录，无法删除"})
		return
	}

	database.DB.Delete(&models.Room{}, id)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetBookings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	date := c.Query("date")
	roomID := c.Query("room_id")
	status := c.Query("status")

	var bookings []models.Booking
	var total int64

	query := database.DB.Model(&models.Booking{}).Preload("Room")

	if date != "" {
		query = query.Where("date = ?", date)
	}
	if roomID != "" {
		query = query.Where("room_id = ?", roomID)
	}
	if status != "" {
		statusInt, _ := strconv.Atoi(status)
		query = query.Where("status = ?", statusInt)
	}

	query.Count(&total)
	query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&bookings)

	c.JSON(http.StatusOK, gin.H{
		"list":  bookings,
		"total": total,
		"page":  page,
	})
}

type CancelRequest struct {
	CancelReason string `json:"cancel_reason" binding:"required"`
}

func CancelBooking(c *gin.Context) {
	id := c.Param("id")
	var booking models.Booking
	if err := database.DB.First(&booking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "预订不存在"})
		return
	}

	var req CancelRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供取消原因"})
		return
	}

	database.DB.Model(&booking).Updates(map[string]interface{}{
		"status":        2,
		"cancel_reason": req.CancelReason,
	})

	c.JSON(http.StatusOK, gin.H{"message": "取消成功"})
}

func GetStats(c *gin.Context) {
	today := time.Now().Format("2006-01-02")

	var todayCount int64
	database.DB.Model(&models.Booking{}).Where("date = ?", today).Count(&todayCount)

	var activeRooms int64
	database.DB.Model(&models.Room{}).Where("status = 1").Count(&activeRooms)

	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	monday := now.AddDate(0, 0, -weekday+1)

	dailyStats := make([]map[string]interface{}, 7)
	for i := 0; i < 7; i++ {
		day := monday.AddDate(0, 0, i).Format("2006-01-02")
		var count int64
		database.DB.Model(&models.Booking{}).Where("date = ?", day).Count(&count)
		dailyStats[i] = map[string]interface{}{
			"date":  day,
			"count": count,
		}
	}

	var rooms []models.Room
	database.DB.Find(&rooms)

	roomStats := make([]map[string]interface{}, 0)
	for _, room := range rooms {
		var count int64
		database.DB.Model(&models.Booking{}).
			Where("room_id = ? AND date >= ?", room.ID, monday.Format("2006-01-02")).
			Count(&count)
		roomStats = append(roomStats, map[string]interface{}{
			"room_name": room.Name,
			"count":     count,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"today_bookings":  todayCount,
		"active_rooms":    activeRooms,
		"daily_stats":     dailyStats,
		"room_stats":      roomStats,
	})
}
