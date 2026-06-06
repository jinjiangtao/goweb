
package models

import "time"

type Admin struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"uniqueIndex"`
	Password  string    `json:"-"`
	Nickname  string    `json:"nickname"`
	CreatedAt time.Time `json:"created_at"`
}

type Room struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Name      string    `json:"name"`
	Capacity  int       `json:"capacity"`
	Devices   string    `json:"devices"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

type Booking struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	RoomID       uint      `json:"room_id"`
	Room         Room      `json:"room" gorm:"foreignKey:RoomID"`
	Name         string    `json:"name"`
	Phone        string    `json:"phone"`
	Date         string    `json:"date"`
	StartTime    string    `json:"start_time"`
	EndTime      string    `json:"end_time"`
	Purpose      string    `json:"purpose"`
	Status       int       `json:"status"`
	CancelReason string    `json:"cancel_reason"`
	CreatedAt    time.Time `json:"created_at"`
}

