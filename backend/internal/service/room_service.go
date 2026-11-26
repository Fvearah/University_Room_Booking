package service

import (
	"context"
	"errors"
	"strings"
	"time"
	"unicode/utf8"

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

// Detail: USER tidak boleh melihat ruangan INACTIVE (dianggap tidak ditemukan).
func (s *RoomService) Detail(ctx context.Context, id uint64, role string) (*model.Room, error) {
	rm, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if role != "ADMIN" && rm.Status != "ACTIVE" {
		return nil, repository.ErrRoomNotFound
	}
	return rm, nil
}

// validLen memeriksa panjang karakter (bukan byte).
func validLen(v string, min, max int) bool {
	n := utf8.RuneCountInString(v)
	return n >= min && n <= max
}

func (s *RoomService) Create(ctx context.Context, req dto.CreateRoomRequest) (*model.Room, error) {
	status := req.Status
	if status == "" {
		status = "ACTIVE"
	}
	if status != "ACTIVE" && status != "INACTIVE" {
		return nil, ErrInvalidRoomInput
	}

	name := strings.TrimSpace(req.Name)
	code := strings.TrimSpace(req.Code)
	location := strings.TrimSpace(req.Location)
	if !validLen(name, 3, 120) || !validLen(code, 2, 30) || !validLen(location, 3, 150) {
		return nil, ErrInvalidRoomInput
	}

	rm := &model.Room{
		Name:        name,
		Code:        code,
		Description: req.Description,
		Capacity:    req.Capacity,
		Location:    location,
		Status:      status,
	}
	return s.repo.Create(ctx, rm)
}

func (s *RoomService) Update(ctx context.Context, id uint64, req dto.UpdateRoomRequest) (*model.Room, error) {
	// Validasi ulang di backend (aturan sama seperti Create)
	if req.Name != nil {
		v := strings.TrimSpace(*req.Name)
		if !validLen(v, 3, 120) {
			return nil, ErrInvalidRoomInput
		}
		req.Name = &v
	}
	if req.Code != nil {
		v := strings.TrimSpace(*req.Code)
		if !validLen(v, 2, 30) {
			return nil, ErrInvalidRoomInput
		}
		req.Code = &v
	}
	if req.Location != nil {
		v := strings.TrimSpace(*req.Location)
		if !validLen(v, 3, 150) {
			return nil, ErrInvalidRoomInput
		}
		req.Location = &v
	}
	if req.Capacity != nil && *req.Capacity == 0 {
		return nil, ErrInvalidRoomInput
	}
	if req.Status != nil && *req.Status != "ACTIVE" && *req.Status != "INACTIVE" {
		return nil, ErrInvalidRoomInput
	}
	return s.repo.Update(ctx, id, req)
}

func (s *RoomService) GetAvailability(ctx context.Context, roomID uint64, dateStr string, role string) (*RoomAvailabilityResponse, error) {
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
	if role != "ADMIN" && rm.Status != "ACTIVE" {
		return nil, repository.ErrRoomNotFound
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
