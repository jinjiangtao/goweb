package ws

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type MessageType string

const (
	MsgDraw      MessageType = "draw"
	MsgCursor    MessageType = "cursor"
	MsgUserJoin  MessageType = "user_join"
	MsgUserLeave MessageType = "user_leave"
	MsgUsers     MessageType = "users"
	MsgClear     MessageType = "clear"
	MsgSync      MessageType = "sync"
	MsgUndo      MessageType = "undo"
	MsgRedo      MessageType = "redo"
	MsgPing      MessageType = "ping"
	MsgPong      MessageType = "pong"
)

type Message struct {
	Type      MessageType `json:"type"`
	Data      interface{} `json:"data"`
	UserID    string      `json:"userId,omitempty"`
	Timestamp int64       `json:"timestamp,omitempty"`
}

type User struct {
	ID     string          `json:"id"`
	Name   string          `json:"name"`
	Avatar string          `json:"avatar"`
	Color  string          `json:"color"`
	Conn   *websocket.Conn `json:"-"`
	Send   chan Message    `json:"-"`
	Cursor *Point          `json:"cursor,omitempty"`
}

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Hub struct {
	boardID    string
	users      map[string]*User
	mu         sync.RWMutex
	Broadcast  chan Message
	Register   chan *User
	Unregister chan *User
}

var hubs = make(map[string]*Hub)
var hubsMu sync.RWMutex

func GetHub(boardID string) *Hub {
	hubsMu.RLock()
	hub, exists := hubs[boardID]
	hubsMu.RUnlock()

	if !exists {
		hubsMu.Lock()
		defer hubsMu.Unlock()
		hub = &Hub{
			boardID:    boardID,
			users:      make(map[string]*User),
			Broadcast:  make(chan Message, 1000),
			Register:   make(chan *User, 100),
			Unregister: make(chan *User, 100),
		}
		hubs[boardID] = hub
		go hub.run()
	}
	return hub
}

func (h *Hub) run() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case user := <-h.Register:
			h.mu.Lock()
			h.users[user.ID] = user
			h.mu.Unlock()
			log.Printf("User %s joined board %s", user.ID, h.boardID)

			h.broadcastUserList()
			h.sendToUser(user.ID, Message{
				Type:      MsgUserJoin,
				Data:      user,
				Timestamp: time.Now().UnixMilli(),
			})

		case user := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.users[user.ID]; ok {
				delete(h.users, user.ID)
				close(user.Send)
				log.Printf("User %s left board %s", user.ID, h.boardID)
			}
			h.mu.Unlock()

			h.broadcastUserList()
			h.Broadcast <- Message{
				Type:      MsgUserLeave,
				Data:      user,
				Timestamp: time.Now().UnixMilli(),
			}

		case message := <-h.Broadcast:
			h.mu.RLock()
			for _, user := range h.users {
				if message.UserID != user.ID {
					select {
					case user.Send <- message:
					default:
						close(user.Send)
						delete(h.users, user.ID)
					}
				}
			}
			h.mu.RUnlock()

		case <-ticker.C:
			h.Broadcast <- Message{
				Type:      MsgPing,
				Timestamp: time.Now().UnixMilli(),
			}
		}
	}
}

func (h *Hub) broadcastUserList() {
	h.mu.RLock()
	users := make([]User, 0, len(h.users))
	for _, u := range h.users {
		users = append(users, *u)
	}
	h.mu.RUnlock()

	h.Broadcast <- Message{
		Type:      MsgUsers,
		Data:      users,
		Timestamp: time.Now().UnixMilli(),
	}
}

func (h *Hub) sendToUser(userID string, message Message) {
	h.mu.RLock()
	user, exists := h.users[userID]
	h.mu.RUnlock()

	if exists {
		user.Send <- message
	}
}

func (h *Hub) Send(msg Message) {
	h.Broadcast <- msg
}

func (h *Hub) GetUserCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.users)
}

func NewUser(name, color string, conn *websocket.Conn) *User {
	return &User{
		ID:     uuid.New().String(),
		Name:   name,
		Avatar: "",
		Color:  color,
		Conn:   conn,
		Send:   make(chan Message, 256),
	}
}

func (u *User) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister <- u
		u.Conn.Close()
	}()

	u.Conn.SetReadLimit(1024 * 1024 * 10)
	u.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	u.Conn.SetPongHandler(func(string) error {
		u.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg Message
		err := u.Conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		msg.UserID = u.ID
		msg.Timestamp = time.Now().UnixMilli()

		if msg.Type == MsgPing {
			u.Send <- Message{
				Type:      MsgPong,
				Timestamp: time.Now().UnixMilli(),
			}
			continue
		}

		if msg.Type == MsgCursor {
			if cursorData, ok := msg.Data.(map[string]interface{}); ok {
				u.Cursor = &Point{
					X: cursorData["x"].(float64),
					Y: cursorData["y"].(float64),
				}
			}
		}

		hub.Send(msg)
	}
}

func (u *User) WritePump() {
	defer u.Conn.Close()

	for {
		select {
		case message, ok := <-u.Send:
			u.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				u.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			data, err := json.Marshal(message)
			if err != nil {
				log.Printf("JSON marshal error: %v", err)
				continue
			}

			err = u.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("Write error: %v", err)
				return
			}
		}
	}
}
