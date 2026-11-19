package service

import (
	"context"
	"errors"
	"time"

	"university-room-booking/internal/dto"
	"university-room-booking/internal/model"
	"university-room-booking/internal/repository"
)

var ErrInvalidRoomInput = errors.New("input data ruangan tidak valid")

type RoomAvailabilityResponse struct {
	RoomID      uint64                  `json:"room_id"`
	RoomName    string                  `json:"room_name"`
	RoomStatus  string                  `json:"room_status"`
	Date        string                  `json:"date"`
	BookedSlots []repository.BookedSlot `json:"booked_slots"`
}

type RoomService struct {
	repo *repository.RoomRepository
}

func NewRoomService(repo *repository.RoomRepository) *RoomService {
	return &RoomService{repo: repo}
}

func (s *RoomService) List(ctx context.Context, role string) ([]model.Room, error) {
	includeInactive := role == "ADMIN"
	return s.repo.FindAll(ctx, includeInactive)
}

func (s *RoomService) Detail(ctx context.Context, id uint64) (*model.Room, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *RoomService) Create(ctx context.Context, req dto.CreateRoomRequest) (*model.Room, error) {
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}
	if status != "ACTIVE" && status != "INACTIVE" {
		return nil, ErrInvalidRoomInput
	}

	rm := &model.Room{
		Name:        req.Name,
		Code:        req.Code,
		Description: req.Description,
		Capacity:    req.Capacity,
		Location:    req.Location,
		Status:      status,
	}
	return s.repo.Create(ctx, rm)
}

func (s *RoomService) Update(ctx context.Context, id uint64, req dto.UpdateRoomRequest) (*model.Room, error) {
	if req.Capacity != nil && *req.Capacity == 0 {
		return nil, ErrInvalidRoomInput
	}
	if req.Status != nil && *req.Status != "ACTIVE" && *req.Status != "INACTIVE" {
		return nil, ErrInvalidRoomInput
	}
	return s.repo.Update(ctx, id, req)
}

func (s *RoomService) GetAvailability(ctx context.Context, roomID uint64, dateStr string) (*RoomAvailabilityResponse, error) {
	loc := time.Local
	if dateStr == "" {
		dateStr = time.Now().In(loc).Format("2006-01-02")
	} else if _, err := time.ParseInLocation("2006-01-02", dateStr, loc); err != nil {
		return nil, ErrInvalidDate
	}

	rm, err := s.repo.FindByID(ctx, roomID)
	if err != nil {
		return nil, err
	}

	slots, err := s.repo.GetBookedSlots(ctx, roomID, dateStr)
	if err != nil {
		return nil, err
	}

	return &RoomAvailabilityResponse{
		RoomID:      rm.ID,
		RoomName:    rm.Name,
		RoomStatus:  rm.Status,
		Date:        dateStr,
		BookedSlots: slots,
	}, nil
}
