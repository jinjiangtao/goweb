package websocket

import (
	"encoding/json"
	"goim/server/cache"
	"goim/server/model"
	"goim/server/storage"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Hub    *Hub
	Conn   *websocket.Conn
	Send   chan []byte
	UserID string
}

type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan *model.WSMessage
	Register   chan *Client
	Unregister chan *Client
}

var hub = &Hub{
	Clients:    make(map[string]*Client),
	Broadcast:  make(chan *model.WSMessage),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client.UserID] = client
			cache.SetOnline(client.UserID)
			log.Printf("User %s connected", client.UserID)
		case client := <-h.Unregister:
			if _, ok := h.Clients[client.UserID]; ok {
				delete(h.Clients, client.UserID)
				close(client.Send)
				cache.SetOffline(client.UserID)
				log.Printf("User %s disconnected", client.UserID)
			}
		case message := <-h.Broadcast:
			h.SendMessage(message)
		}
	}
}

func (h *Hub) SendMessage(msg *model.WSMessage) {
	log.Printf("SendMessage: from=%s, to=%s, toType=%d", msg.From, msg.To, msg.ToType)

	if msg.ToType == 0 {
		if client, ok := h.Clients[msg.To]; ok {
			log.Printf("Sending to user %s, client found", msg.To)
			select {
			case client.Send <- serializeMessage(msg):
				log.Printf("Message sent to %s", msg.To)
			default:
				log.Printf("Failed to send to %s, closing client", msg.To)
				close(client.Send)
				delete(h.Clients, client.UserID)
			}
		} else {
			log.Printf("Client %s not found in hub", msg.To)
		}
	} else {
		members, err := storage.GetGroupMembers(msg.To)
		if err != nil {
			log.Printf("Failed to get group members: %v", err)
			return
		}
		for _, member := range members {
			if member.UserID != msg.From {
				if client, ok := h.Clients[member.UserID]; ok {
					select {
					case client.Send <- serializeMessage(msg):
					default:
						close(client.Send)
						delete(h.Clients, client.UserID)
					}
				}
			}
		}
	}
}

func serializeMessage(msg *model.WSMessage) []byte {
	data, _ := json.Marshal(msg)
	return data
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	userID := r.URL.Query().Get("userID")
	if userID == "" {
		http.Error(w, "userID is required", http.StatusBadRequest)
		return
	}

	log.Printf("WebSocket connection request from user: %s", userID)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	client := &Client{
		Hub:    hub,
		Conn:   conn,
		Send:   make(chan []byte, 256),
		UserID: userID,
	}

	client.Hub.Register <- client

	go client.WritePump()
	go client.ReadPump()
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(512)
	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))

	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error reading message: %v", err)
			}
			break
		}

		log.Printf("Received message from %s: %s", c.UserID, string(message))

		var wsMsg model.WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("Error parsing message: %v", err)
			continue
		}

		wsMsg.From = c.UserID

		if wsMsg.Type == "ping" {
			c.Send <- serializeMessage(&model.WSMessage{Type: "pong"})
			continue
		}

		handleMessage(&wsMsg)
	}
}

func handleMessage(msg *model.WSMessage) {
	log.Printf("handleMessage: from=%s, to=%s, toType=%d, content=%s", msg.From, msg.To, msg.ToType, msg.Content)

	dbMsg := &model.Message{
		ID:           uuid.New().String(),
		SenderID:     msg.From,
		ReceiverID:   msg.To,
		ReceiverType: msg.ToType,
		Content:      msg.Content,
		Type:         msg.MsgType,
		Status:       0,
		CreatedAt:    time.Now(),
	}

	err := storage.CreateMessage(dbMsg)
	if err != nil {
		log.Printf("Failed to save message: %v", err)
		return
	}

	if !cache.IsOnline(msg.To) && msg.ToType == 0 {
		cache.IncrUnreadCount(msg.To, msg.From)
	}

	msg.ID = dbMsg.ID
	msg.Timestamp = dbMsg.CreatedAt.Unix()
	hub.Broadcast <- msg
	log.Printf("Message broadcasted: id=%s, to=%s", msg.ID, msg.To)
}

func (c *Client) WritePump() {
	ticker := time.NewTicker(60 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func GetHub() *Hub {
	return hub
}
