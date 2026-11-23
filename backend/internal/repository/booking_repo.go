package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"university-room-booking/internal/model"
)

var (
	ErrRoomNotFound    = errors.New("room not found")
	ErrRoomInactive    = errors.New("room is inactive")
	ErrRoomConflict    = errors.New("room schedule conflict")
	ErrBookingNotFound = errors.New("booking not found")
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) CreateBookingTx(ctx context.Context, b *model.Booking) (*model.Booking, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// 1. Kunci baris ruangan untuk mencegah race condition
	var roomID uint64
	var roomStatus string
	err = tx.QueryRowContext(ctx, `SELECT id, status FROM rooms WHERE id = ? FOR UPDATE`, b.RoomID).Scan(&roomID, &roomStatus)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRoomNotFound
	}
	if err != nil {
		return nil, err
	}
	if roomStatus != "ACTIVE" {
		err = ErrRoomInactive
		return nil, err
	}

	// 2. Cek bentrok jadwal pada booking CONFIRMED
	var conflictID uint64
	conflictQuery := `
		SELECT id FROM bookings
		WHERE room_id = ?
		  AND booking_date = ?
		  AND status = 'CONFIRMED'
		  AND start_time < ?
		  AND end_time > ?
		LIMIT 1`
	err = tx.QueryRowContext(ctx, conflictQuery, b.RoomID, b.BookingDate, b.EndTime, b.StartTime).Scan(&conflictID)
	if err == nil {
		err = ErrRoomConflict
		return nil, err
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return nil, err
	}

	// 3. Simpan booking baru dengan status CONFIRMED
	insertQuery := `
		INSERT INTO bookings (user_id, room_id, title, purpose, booking_date, start_time, end_time, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, 'CONFIRMED', NOW(), NOW())`
	res, err := tx.ExecContext(ctx, insertQuery, b.UserID, b.RoomID, b.Title, b.Purpose, b.BookingDate, b.StartTime, b.EndTime)
	if err != nil {
		return nil, err
	}

	lastID, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	b.ID = uint64(lastID)
	b.Status = "CONFIRMED"
	return b, nil
}

func (r *BookingRepository) FindAll(ctx context.Context, userID *uint64, dateFilter, roomIDFilter, statusFilter string) ([]model.Booking, error) {
	query := `
		SELECT b.id, b.user_id, u.name, b.room_id, rm.name, rm.code,
		       b.title, b.purpose, b.booking_date, b.start_time, b.end_time, b.status, b.created_at
		FROM bookings b
		JOIN users u ON b.user_id = u.id
		JOIN rooms rm ON b.room_id = rm.id`

	var conditions []string
	var args []interface{}

	if userID != nil {
		conditions = append(conditions, "b.user_id = ?")
		args = append(args, *userID)
	}
	if dateFilter != "" {
		conditions = append(conditions, "b.booking_date = ?")
		args = append(args, dateFilter)
	}
	if roomIDFilter != "" {
		conditions = append(conditions, "b.room_id = ?")
		args = append(args, roomIDFilter)
	}
	if statusFilter != "" && statusFilter != "COMPLETED" {
		conditions = append(conditions, "b.status = ?")
		args = append(args, statusFilter)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY b.booking_date DESC, b.start_time ASC"

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bookings := make([]model.Booking, 0)
	for rows.Next() {
		var b model.Booking
		var bDate time.Time
		if err := rows.Scan(
			&b.ID, &b.UserID, &b.UserName, &b.RoomID, &b.RoomName, &b.RoomCode,
			&b.Title, &b.Purpose, &bDate, &b.StartTime, &b.EndTime, &b.Status, &b.CreatedAt,
		); err != nil {
			return nil, err
		}
		b.BookingDate = bDate.Format("2006-01-02")
		bookings = append(bookings, b)
	}
	return bookings, rows.Err()
}

func (r *BookingRepository) FindByID(ctx context.Context, id uint64) (*model.Booking, error) {
	query := `
		SELECT b.id, b.user_id, u.name, b.room_id, rm.name, rm.code,
		       b.title, b.purpose, b.booking_date, b.start_time, b.end_time, b.status, b.created_at
		FROM bookings b
		JOIN users u ON b.user_id = u.id
		JOIN rooms rm ON b.room_id = rm.id
		WHERE b.id = ?`

	var b model.Booking
	var bDate time.Time
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&b.ID, &b.UserID, &b.UserName, &b.RoomID, &b.RoomName, &b.RoomCode,
		&b.Title, &b.Purpose, &bDate, &b.StartTime, &b.EndTime, &b.Status, &b.CreatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrBookingNotFound
	}
	if err != nil {
		return nil, err
	}
	b.BookingDate = bDate.Format("2006-01-02")
	return &b, nil
}

func (r *BookingRepository) CancelBooking(ctx context.Context, id uint64) error {
	res, err := r.db.ExecContext(ctx, `UPDATE bookings SET status = 'CANCELLED', updated_at = NOW() WHERE id = ? AND status = 'CONFIRMED'`, id)
	if err != nil {
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return fmt.Errorf("booking tidak dapat dibatalkan")
	}
	return nil
}
