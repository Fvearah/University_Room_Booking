package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"university-room-booking/internal/dto"
	"university-room-booking/internal/model"
	"university-room-booking/internal/repository"
)

var (
	ErrInvalidDate     = errors.New("tanggal pemakaian tidak valid atau berada di masa lalu")
	ErrInvalidTime     = errors.New("jam mulai harus lebih awal daripada jam selesai")
	ErrForbiddenAccess = errors.New("tidak memiliki hak akses terhadap booking ini")
	ErrCannotCancel    = errors.New("booking yang sudah dimulai, selesai, atau dibatalkan tidak dapat dibatalkan")
)

type BookingService struct {
	repo *repository.BookingRepository
}

func NewBookingService(repo *repository.BookingRepository) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) Create(ctx context.Context, userID uint64, req dto.CreateBookingRequest) (*model.Booking, error) {
	loc := time.Local
	now := time.Now().In(loc)

	parsedDate, err := time.ParseInLocation("2006-01-02", req.BookingDate, loc)
	if err != nil {
		return nil, ErrInvalidDate
	}

	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	if parsedDate.Before(today) {
		return nil, ErrInvalidDate
	}

	startNorm, startClock, err := normalizeTime(req.StartTime)
	if err != nil {
		return nil, ErrInvalidTime
	}
	endNorm, endClock, err := normalizeTime(req.EndTime)
	if err != nil {
		return nil, ErrInvalidTime
	}
	if !startClock.Before(endClock) {
		return nil, ErrInvalidTime
	}

	// Tolak jika booking hari ini namun jam mulai sudah terlewati
	startFull := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), startClock.Hour(), startClock.Minute(), 0, 0, loc)
	if !startFull.After(now) {
		return nil, ErrInvalidDate
	}

	booking := &model.Booking{
		UserID:      userID,
		RoomID:      req.RoomID,
		Title:       req.Title,
		Purpose:     req.Purpose,
		BookingDate: req.BookingDate,
		StartTime:   startNorm,
		EndTime:     endNorm,
	}

	return s.repo.CreateBookingTx(ctx, booking)
}

func (s *BookingService) List(ctx context.Context, userID uint64, role, dateFilter, roomIDFilter, statusFilter string) ([]model.Booking, error) {
	var ownerFilter *uint64
	if role != "ADMIN" {
		ownerFilter = &userID
	}

	items, err := s.repo.FindAll(ctx, ownerFilter, dateFilter, roomIDFilter, statusFilter)
	if err != nil {
		return nil, err
	}

	now := time.Now().In(time.Local)
	filtered := make([]model.Booking, 0, len(items))
	for _, b := range items {
		b.Status = deriveDisplayStatus(b, now)
		if (statusFilter == "COMPLETED" || statusFilter == "CONFIRMED") && b.Status != statusFilter {
			continue

		}
		filtered = append(filtered, b)
	}
	return filtered, nil
}

func (s *BookingService) GetDetail(ctx context.Context, bookingID, userID uint64, role string) (*model.Booking, error) {
	b, err := s.repo.FindByID(ctx, bookingID)
	if err != nil {
		return nil, err
	}
	if role != "ADMIN" && b.UserID != userID {
		return nil, ErrForbiddenAccess
	}
	b.Status = deriveDisplayStatus(*b, time.Now().In(time.Local))
	return b, nil
}

func (s *BookingService) Cancel(ctx context.Context, bookingID, userID uint64) error {
	b, err := s.repo.FindByID(ctx, bookingID)
	if err != nil {
		return err
	}
	if b.UserID != userID {
		return ErrForbiddenAccess
	}
	if b.Status != "CONFIRMED" {
		return ErrCannotCancel
	}

	loc := time.Local
	startDT, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %s", b.BookingDate, b.StartTime), loc)
	if err != nil || !time.Now().In(loc).Before(startDT) {
		return ErrCannotCancel
	}

	return s.repo.CancelBooking(ctx, bookingID)
}

func normalizeTime(raw string) (string, time.Time, error) {
	for _, layout := range []string{"15:04", "15:04:05"} {
		if t, err := time.Parse(layout, raw); err == nil {
			return t.Format("15:04:05"), t, nil
		}
	}
	return "", time.Time{}, fmt.Errorf("format waktu tidak valid")
}

func deriveDisplayStatus(b model.Booking, now time.Time) string {
	if b.Status != "CONFIRMED" {
		return b.Status
	}
	endDT, err := time.ParseInLocation("2006-01-02 15:04:05", fmt.Sprintf("%s %s", b.BookingDate, b.EndTime), time.Local)
	if err == nil && now.After(endDT) {
		return "COMPLETED"
	}
	return "CONFIRMED"
}
