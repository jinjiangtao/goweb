package handlers

import (
	"huiyishi-server/database"
	"huiyishi-server/models"
	"huiyishi-server/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func AdminLogin(c *gin.Context) {
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

func UserLogin(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	if user.Status != 1 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "账号已被禁用"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "用户名或密码错误"})
		return
	}

	token, err := utils.GenerateUserToken(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "生成token失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":        user.ID,
			"username":  user.Username,
			"real_name": user.RealName,
			"phone":     user.Phone,
		},
	})
}

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RealName string `json:"real_name" binding:"required"`
	Phone    string `json:"phone" binding:"required,len=11"`
}

func UserRegister(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var count int64
	database.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := models.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		RealName:  req.RealName,
		Phone:     req.Phone,
		Status:    1,
		CreatedAt: time.Now(),
	}

	database.DB.Create(&user)
	c.JSON(http.StatusOK, user)
}

func GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("user_id")
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":        user.ID,
		"username":  user.Username,
		"real_name": user.RealName,
		"phone":     user.Phone,
	})
}

func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")

	var users []models.User
	var total int64

	query := database.DB.Model(&models.User{})
	if search != "" {
		query = query.Where("username LIKE ? OR real_name LIKE ? OR phone LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	query.Count(&total)
	query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"list":  users,
		"total": total,
		"page":  page,
	})
}

type CreateUserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	RealName string `json:"real_name" binding:"required"`
	Phone    string `json:"phone" binding:"required,len=11"`
}

func CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var count int64
	database.DB.Model(&models.User{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户名已存在"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	user := models.User{
		Username:  req.Username,
		Password:  string(hashedPassword),
		RealName:  req.RealName,
		Phone:     req.Phone,
		Status:    1,
		CreatedAt: time.Now(),
	}

	database.DB.Create(&user)
	c.JSON(http.StatusOK, user)
}

type UpdateUserRequest struct {
	RealName string `json:"real_name"`
	Phone    string `json:"phone"`
}

func UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	database.DB.Model(&user).Updates(map[string]interface{}{
		"real_name": req.RealName,
		"phone":     req.Phone,
	})

	c.JSON(http.StatusOK, user)
}

func ToggleUserStatus(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	newStatus := 1
	if user.Status == 1 {
		newStatus = 0
	}

	database.DB.Model(&user).Updates(map[string]interface{}{
		"status": newStatus,
	})

	c.JSON(http.StatusOK, gin.H{"message": "状态更新成功"})
}

type ResetPasswordRequest struct {
	Password string `json:"password" binding:"required"`
}

func ResetUserPassword(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	database.DB.Model(&user).Updates(map[string]interface{}{
		"password": string(hashedPassword),
	})

	c.JSON(http.StatusOK, gin.H{"message": "密码重置成功"})
}

func DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := database.DB.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	database.DB.Delete(&user)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func GetRooms(c *gin.Context) {
	search := c.Query("search")

	var rooms []models.Room
	query := database.DB.Order("id desc")

	if search != "" {query = query.Where("name LIKE ?", "%"+search+"%")
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

	query := database.DB.Model(&models.Booking{}).Preload("Room").Preload("User")

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
		"today_bookings": todayCount,
		"active_rooms":   activeRooms,
		"daily_stats":    dailyStats,
		"room_stats":     roomStats,
	})
}

func PublicGetRooms(c *gin.Context) {
	var rooms []models.Room
	database.DB.Where("status = ?", 1).Find(&rooms)
	c.JSON(http.StatusOK, rooms)
}

func PublicGetRoomAvailable(c *gin.Context) {
	roomID := c.Param("id")
	date := c.Query("date")

	if roomID == "" || date == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var bookings []models.Booking
	database.DB.Where("room_id = ? AND date = ? AND status = ?", roomID, date, 1).Find(&bookings)

	var occupiedSlots []map[string]string
	for _, booking := range bookings {
		occupiedSlots = append(occupiedSlots, map[string]string{
			"start_time": booking.StartTime,
			"end_time":   booking.EndTime,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"occupied": occupiedSlots,
	})
}

type PublicBookingRequest struct {
	RoomID    uint   `json:"room_id" binding:"required"`
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	Date      string `json:"date" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
	Purpose   string `json:"purpose"`
}

func PublicCreateBooking(c *gin.Context) {
	var req PublicBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	var count int64
	database.DB.Model(&models.Booking{}).
		Where("room_id = ? AND date = ? AND status = ?", req.RoomID, req.Date, 1).
		Where("(? < end_time AND ? > start_time)", req.StartTime, req.EndTime).
		Count(&count)

	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "该时间段已被预订，请更换时间"})
		return
	}

	var userID *uint
	var name, phone string

	if uid, exists := c.Get("user_id"); exists {
		var user models.User
		if database.DB.First(&user, uid).Error == nil {
			uidUint := uid.(uint)
			userID = &uidUint
			name = user.RealName
			phone = user.Phone
		}
	}

	if name == "" {
		name = req.Name
	}
	if phone == "" {
		phone = req.Phone
	}

	if name == "" || phone == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供姓名和手机号"})
		return
	}

	booking := models.Booking{
		RoomID:    req.RoomID,
		UserID:    userID,
		Name:      name,
		Phone:     phone,
		Date:      req.Date,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Purpose:   req.Purpose,
		Status:    1,
		CreatedAt: time.Now(),
	}

	database.DB.Create(&booking)
	c.JSON(http.StatusOK, booking)
}

func PublicGetMyBookings(c *gin.Context) {
	var bookings []models.Booking

	if uid, exists := c.Get("user_id"); exists {
		database.DB.Where("user_id = ?", uid).Preload("Room").Order("id desc").Find(&bookings)
	} else {
		phone := c.Query("phone")
		if phone == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "请登录或提供手机号"})
			return
		}
		database.DB.Where("phone = ?", phone).Preload("Room").Order("id desc").Find(&bookings)
	}

	c.JSON(http.StatusOK, bookings)
}

func PublicCancelBooking(c *gin.Context) {
	id := c.Param("id")

	var booking models.Booking
	if err := database.DB.First(&booking, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "预订不存在"})
		return
	}

	hasPermission := false

	if _, exists := c.Get("admin_id"); exists {
		hasPermission = true
	} else if uid, exists := c.Get("user_id"); exists {
		if booking.UserID != nil && *booking.UserID == uid.(uint) {
			hasPermission = true
		}
	} else {
		phone := c.Query("phone")
		if phone != "" && booking.Phone == phone {
			hasPermission = true
		}
	}

	if !hasPermission {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "无权取消此预订"})
		return
	}

	if booking.Status != 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "预订状态不允许取消"})
		return
	}

	database.DB.Model(&booking).Updates(map[string]interface{}{
		"status": 2,
	})

	c.JSON(http.StatusOK, gin.H{"message": "取消成功"})
}
