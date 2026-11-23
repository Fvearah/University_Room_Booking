// internal/dto/dto.go
package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type CreateRoomRequest struct {
	Name        string `json:"name" binding:"required,min=3,max=120"`
	Code        string `json:"code" binding:"required,min=2,max=30"`
	Description string `json:"description"`
	Capacity    uint   `json:"capacity" binding:"required,gt=0"`
	Location    string `json:"location" binding:"required,min=3,max=150"`
	Status      string `json:"status"`
}

type UpdateRoomRequest struct {
	Name        *string `json:"name"`
	Code        *string `json:"code"`
	Description *string `json:"description"`
	Capacity    *uint   `json:"capacity"`
	Location    *string `json:"location"`
	Status      *string `json:"status"`
}

type CreateBookingRequest struct {
	RoomID      uint64 `json:"room_id" binding:"required"`
	Title       string `json:"title" binding:"required,min=5,max=150"`
	Purpose     string `json:"purpose" binding:"required,min=10,max=2000"`
	BookingDate string `json:"booking_date" binding:"required"`
	StartTime   string `json:"start_time" binding:"required"`
	EndTime     string `json:"end_time" binding:"required"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}
