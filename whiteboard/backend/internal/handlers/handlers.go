package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"whiteboard/internal/models"
	"whiteboard/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
	ReadBufferSize:  1024 * 1024,
	WriteBufferSize: 1024 * 1024,
}

func HandleWebSocket(c *gin.Context) {
	boardID := c.Param("id")
	userName := c.DefaultQuery("name", "匿名用户")
	userColor := c.DefaultQuery("color", "#409eff")

	if boardID == "" {
		c.JSON(400, gin.H{"error": "Board ID is required"})
		return
	}

	var board models.Board
	if err := models.DB.Where("id = ?", boardID).First(&board).Error; err != nil {
		c.JSON(404, gin.H{"error": "Board not found"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	hub := ws.GetHub(boardID)
	user := ws.NewUser(userName, userColor, conn)

	userInfo := map[string]interface{}{
		"id":     user.ID,
		"name":   user.Name,
		"color":  user.Color,
		"avatar": user.Avatar,
	}
	user.Send <- ws.Message{
		Type:      ws.MsgUserInfo,
		Data:      userInfo,
		Timestamp: time.Now().UnixMilli(),
	}

	hub.Register <- user

	syncData := getBoardSyncData(boardID)
	user.Send <- ws.Message{
		Type:      ws.MsgSync,
		Data:      syncData,
		Timestamp: time.Now().UnixMilli(),
	}

	go user.WritePump()
	user.ReadPump(hub)
}

func getBoardSyncData(boardID string) gin.H {
	var elements []models.Element
	models.DB.Where("board_id = ?", boardID).Order("timestamp asc").Find(&elements)

	result := make([]map[string]interface{}, 0, len(elements))
	for _, elem := range elements {
		elem.Decode()
		elemMap := map[string]interface{}{
			"id":         elem.ID,
			"type":       elem.Type,
			"userId":     elem.UserID,
			"points":     elem.PointsArr,
			"style":      elem.StyleObj,
			"timestamp":  elem.Timestamp,
			"text":       elem.Text,
			"fontSize":   elem.FontSize,
			"fontFamily": elem.FontFamily,
			"imageData":  elem.ImageData,
			"x":          elem.X,
			"y":          elem.Y,
			"width":      elem.Width,
			"height":     elem.Height,
		}
		result = append(result, elemMap)
	}

	return gin.H{
		"elements": result,
	}
}

func ListBoards(c *gin.Context) {
	var boards []models.Board
	if err := models.DB.Order("updated_at desc").Find(&boards).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": boards})
}

func GetBoard(c *gin.Context) {
	id := c.Param("id")
	var board models.Board
	if err := models.DB.Where("id = ?", id).First(&board).Error; err != nil {
		c.JSON(404, gin.H{"error": "Board not found"})
		return
	}
	c.JSON(200, gin.H{"data": board})
}

func CreateBoard(c *gin.Context) {
	var req struct {
		Name       string `json:"name" binding:"required"`
		Background string `json:"background"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	board := models.Board{
		ID:         uuid.New().String(),
		Name:       req.Name,
		Background: req.Background,
	}

	if board.Background == "" {
		board.Background = "#ffffff"
	}

	if err := models.DB.Create(&board).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	logOperation(board.ID, "system", "create_board", req.Name)

	c.JSON(201, gin.H{"data": board})
}

func UpdateBoard(c *gin.Context) {
	id := c.Param("id")
	var req struct {
		Name       string `json:"name"`
		Background string `json:"background"`
		Thumbnail  string `json:"thumbnail"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	var board models.Board
	if err := models.DB.Where("id = ?", id).First(&board).Error; err != nil {
		c.JSON(404, gin.H{"error": "Board not found"})
		return
	}

	updates := make(map[string]interface{})
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Background != "" {
		updates["background"] = req.Background
	}
	if req.Thumbnail != "" {
		updates["thumbnail"] = req.Thumbnail
	}
	updates["updated_at"] = time.Now()

	if err := models.DB.Model(&board).Updates(updates).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	logOperation(id, "system", "update_board", "")

	c.JSON(200, gin.H{"data": board})
}

func DeleteBoard(c *gin.Context) {
	id := c.Param("id")
	if err := models.DB.Where("id = ?", id).Delete(&models.Board{}).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"message": "Board deleted successfully"})
}

func SaveElements(c *gin.Context) {
	boardID := c.Param("id")
	var req struct {
		Elements []map[string]interface{} `json:"elements" binding:"required"`
		UserID   string                   `json:"userId"`
		UserName string                   `json:"userName"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := models.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Where("board_id = ?", boardID).Delete(&models.Element{}).Error; err != nil {
		tx.Rollback()
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	for _, elemData := range req.Elements {
		elem := models.Element{
			BoardID:   boardID,
			Timestamp: time.Now().UnixMilli(),
		}

		if id, ok := elemData["id"].(string); ok {
			elem.ID = id
		}
		if typ, ok := elemData["type"].(string); ok {
			elem.Type = typ
		}
		if uid, ok := elemData["userId"].(string); ok {
			elem.UserID = uid
		}
		if text, ok := elemData["text"].(string); ok {
			elem.Text = text
		}
		if fs, ok := elemData["fontSize"].(float64); ok {
			elem.FontSize = int(fs)
		}
		if ff, ok := elemData["fontFamily"].(string); ok {
			elem.FontFamily = ff
		}
		if img, ok := elemData["imageData"].(string); ok {
			elem.ImageData = img
		}
		if x, ok := elemData["x"].(float64); ok {
			elem.X = x
		}
		if y, ok := elemData["y"].(float64); ok {
			elem.Y = y
		}
		if w, ok := elemData["width"].(float64); ok {
			elem.Width = w
		}
		if h, ok := elemData["height"].(float64); ok {
			elem.Height = h
		}

		if points, ok := elemData["points"].([]interface{}); ok {
			pointsJSON, _ := json.Marshal(points)
			elem.Points = string(pointsJSON)
			json.Unmarshal(pointsJSON, &elem.PointsArr)
		}
		if style, ok := elemData["style"].(map[string]interface{}); ok {
			styleJSON, _ := json.Marshal(style)
			elem.Style = string(styleJSON)
			json.Unmarshal(styleJSON, &elem.StyleObj)
		}

		if err := tx.Create(&elem).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	logOperation(boardID, req.UserID, "save_elements", "")
	saveHistory(boardID, req.UserID, req.UserName, "save")

	c.JSON(200, gin.H{"message": "Elements saved successfully"})
}

func GetElements(c *gin.Context) {
	boardID := c.Param("id")
	var elements []models.Element
	if err := models.DB.Where("board_id = ?", boardID).Order("timestamp asc").Find(&elements).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	result := make([]map[string]interface{}, 0, len(elements))
	for _, elem := range elements {
		elem.Decode()
		elemMap := map[string]interface{}{
			"id":         elem.ID,
			"type":       elem.Type,
			"userId":     elem.UserID,
			"points":     elem.PointsArr,
			"style":      elem.StyleObj,
			"timestamp":  elem.Timestamp,
			"text":       elem.Text,
			"fontSize":   elem.FontSize,
			"fontFamily": elem.FontFamily,
			"imageData":  elem.ImageData,
			"x":          elem.X,
			"y":          elem.Y,
			"width":      elem.Width,
			"height":     elem.Height,
		}
		result = append(result, elemMap)
	}

	c.JSON(200, gin.H{"data": result})
}

func ClearBoard(c *gin.Context) {
	boardID := c.Param("id")
	var req struct {
		UserID   string `json:"userId"`
		UserName string `json:"userName"`
	}
	c.ShouldBindJSON(&req)

	saveHistory(boardID, req.UserID, req.UserName, "clear")

	if err := models.DB.Where("board_id = ?", boardID).Delete(&models.Element{}).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	logOperation(boardID, req.UserID, "clear_board", "")

	hub := ws.GetHub(boardID)
	hub.Send(ws.Message{
		Type:      ws.MsgClear,
		UserID:    req.UserID,
		Timestamp: time.Now().UnixMilli(),
	})

	c.JSON(200, gin.H{"message": "Board cleared successfully"})
}

func GetHistory(c *gin.Context) {
	boardID := c.Param("id")
	var histories []models.History
	if err := models.DB.Where("board_id = ?", boardID).Order("timestamp desc").Limit(50).Find(&histories).Error; err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"data": histories})
}

func RestoreHistory(c *gin.Context) {
	boardID := c.Param("id")
	historyID := c.Param("historyId")
	var req struct {
		UserID   string `json:"userId"`
		UserName string `json:"userName"`
	}
	c.ShouldBindJSON(&req)

	var history models.History
	if err := models.DB.Where("id = ? AND board_id = ?", historyID, boardID).First(&history).Error; err != nil {
		c.JSON(404, gin.H{"error": "History not found"})
		return
	}

	if history.Snapshot != "" {
		tx := models.DB.Begin()

		if err := tx.Where("board_id = ?", boardID).Delete(&models.Element{}).Error; err != nil {
			tx.Rollback()
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		var elements []map[string]interface{}
		if err := json.Unmarshal([]byte(history.Snapshot), &elements); err == nil {
			for _, elemData := range elements {
				elem := models.Element{
					BoardID:   boardID,
					Timestamp: time.Now().UnixMilli(),
				}
				if id, ok := elemData["id"].(string); ok {
					elem.ID = id
				}
				if typ, ok := elemData["type"].(string); ok {
					elem.Type = typ
				}
				if uid, ok := elemData["userId"].(string); ok {
					elem.UserID = uid
				}
				if text, ok := elemData["text"].(string); ok {
					elem.Text = text
				}
				if points, ok := elemData["points"].([]interface{}); ok {
					pointsJSON, _ := json.Marshal(points)
					elem.Points = string(pointsJSON)
				}
				if style, ok := elemData["style"].(map[string]interface{}); ok {
					styleJSON, _ := json.Marshal(style)
					elem.Style = string(styleJSON)
				}
				tx.Create(&elem)
			}
		}

		tx.Commit()

		hub := ws.GetHub(boardID)
		syncData := getBoardSyncData(boardID)
		hub.Send(ws.Message{
			Type:      ws.MsgSync,
			Data:      syncData,
			UserID:    req.UserID,
			Timestamp: time.Now().UnixMilli(),
		})
	}

	logOperation(boardID, req.UserID, "restore_history", historyID)

	c.JSON(200, gin.H{"message": "History restored successfully"})
}

func saveHistory(boardID, userID, userName, action string) {
	var elements []models.Element
	models.DB.Where("board_id = ?", boardID).Find(&elements)

	snapshotData := make([]map[string]interface{}, 0, len(elements))
	for _, elem := range elements {
		elem.Decode()
		elemMap := map[string]interface{}{
			"id":         elem.ID,
			"type":       elem.Type,
			"userId":     elem.UserID,
			"points":     elem.PointsArr,
			"style":      elem.StyleObj,
			"text":       elem.Text,
			"fontSize":   elem.FontSize,
			"fontFamily": elem.FontFamily,
			"imageData":  elem.ImageData,
		}
		snapshotData = append(snapshotData, elemMap)
	}

	snapshotJSON, _ := json.Marshal(snapshotData)

	history := models.History{
		ID:        uuid.New().String(),
		BoardID:   boardID,
		Action:    action,
		UserID:    userID,
		UserName:  userName,
		Snapshot:  string(snapshotJSON),
		Timestamp: time.Now(),
	}

	models.DB.Create(&history)

	var count int64
	models.DB.Model(&models.History{}).Where("board_id = ?", boardID).Count(&count)
	if count > 100 {
		var oldest models.History
		models.DB.Where("board_id = ?", boardID).Order("timestamp asc").First(&oldest)
		models.DB.Delete(&oldest)
	}
}

func logOperation(boardID, userID, action, data string) {
	log := models.OperationLog{
		ID:        uuid.New().String(),
		BoardID:   boardID,
		UserID:    userID,
		Action:    action,
		Data:      data,
		Timestamp: time.Now(),
	}
	models.DB.Create(&log)
}
