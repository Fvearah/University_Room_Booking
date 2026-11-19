package model

import "time"

type User struct {
	ID           uint64    `json:"id"`
	Name         string    `json:"name"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	Role         string    `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
}

type Room struct {
	ID          uint64    `json:"id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Description string    `json:"description"`
	Capacity    uint      `json:"capacity"`
	Location    string    `json:"location"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Booking struct {
	ID          uint64    `json:"id"`
	UserID      uint64    `json:"user_id"`
	UserName    string    `json:"user_name,omitempty"`
	RoomID      uint64    `json:"room_id"`
	RoomName    string    `json:"room_name,omitempty"`
	RoomCode    string    `json:"room_code,omitempty"`
	Title       string    `json:"title"`
	Purpose     string    `json:"purpose"`
	BookingDate string    `json:"booking_date"`
	StartTime   string    `json:"start_time"`
	EndTime     string    `json:"end_time"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}
