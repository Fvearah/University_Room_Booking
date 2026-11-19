package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"university-room-booking/internal/dto"
	"university-room-booking/internal/model"

	"github.com/go-sql-driver/mysql"
)

var ErrDuplicateRoomCode = errors.New("kode ruangan sudah digunakan")

type BookedSlot struct {
	BookingID uint64 `json:"booking_id"`
	Title     string `json:"title"`
	UserName  string `json:"user_name"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) FindAll(ctx context.Context, includeInactive bool) ([]model.Room, error) {
	query := `SELECT id, name, code, COALESCE(description, ''), capacity, location, status, created_at, updated_at FROM rooms`
	if !includeInactive {
		query += ` WHERE status = 'ACTIVE'`
	}
	query += ` ORDER BY name ASC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rooms := make([]model.Room, 0)
	for rows.Next() {
		var rm model.Room
		if err := rows.Scan(
			&rm.ID, &rm.Name, &rm.Code, &rm.Description, &rm.Capacity,
			&rm.Location, &rm.Status, &rm.CreatedAt, &rm.UpdatedAt,
		); err != nil {
			return nil, err
		}
		rooms = append(rooms, rm)
	}
	return rooms, rows.Err()
}

func (r *RoomRepository) FindByID(ctx context.Context, id uint64) (*model.Room, error) {
	query := `SELECT id, name, code, COALESCE(description, ''), capacity, location, status, created_at, updated_at FROM rooms WHERE id = ? LIMIT 1`
	var rm model.Room
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rm.ID, &rm.Name, &rm.Code, &rm.Description, &rm.Capacity,
		&rm.Location, &rm.Status, &rm.CreatedAt, &rm.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrRoomNotFound
	}
	if err != nil {
		return nil, err
	}
	return &rm, nil
}

func (r *RoomRepository) Create(ctx context.Context, rm *model.Room) (*model.Room, error) {
	query := `
		INSERT INTO rooms (name, code, description, capacity, location, status, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, NOW(), NOW())`
	res, err := r.db.ExecContext(ctx, query, rm.Name, rm.Code, rm.Description, rm.Capacity, rm.Location, rm.Status)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrDuplicateRoomCode
		}
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	return r.FindByID(ctx, uint64(id))
}

func (r *RoomRepository) Update(ctx context.Context, id uint64, req dto.UpdateRoomRequest) (*model.Room, error) {
	var sets []string
	var args []interface{}

	if req.Name != nil {
		sets = append(sets, "name = ?")
		args = append(args, *req.Name)
	}
	if req.Code != nil {
		sets = append(sets, "code = ?")
		args = append(args, *req.Code)
	}
	if req.Description != nil {
		sets = append(sets, "description = ?")
		args = append(args, *req.Description)
	}
	if req.Capacity != nil {
		sets = append(sets, "capacity = ?")
		args = append(args, *req.Capacity)
	}
	if req.Location != nil {
		sets = append(sets, "location = ?")
		args = append(args, *req.Location)
	}
	if req.Status != nil {
		sets = append(sets, "status = ?")
		args = append(args, *req.Status)
	}

	if len(sets) == 0 {
		return r.FindByID(ctx, id)
	}

	sets = append(sets, "updated_at = NOW()")
	args = append(args, id)

	query := "UPDATE rooms SET " + strings.Join(sets, ", ") + " WHERE id = ?"
	res, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil, ErrDuplicateRoomCode
		}
		return nil, err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return r.FindByID(ctx, id)
	}

	return r.FindByID(ctx, id)
}

func (r *RoomRepository) GetBookedSlots(ctx context.Context, roomID uint64, date string) ([]BookedSlot, error) {
	query := `
		SELECT b.id, b.title, u.name, b.start_time, b.end_time
		FROM bookings b
		JOIN users u ON b.user_id = u.id
		WHERE b.room_id = ?
		  AND b.booking_date = ?
		  AND b.status = 'CONFIRMED'
		ORDER BY b.start_time ASC`

	rows, err := r.db.QueryContext(ctx, query, roomID, date)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	slots := make([]BookedSlot, 0)
	for rows.Next() {
		var s BookedSlot
		if err := rows.Scan(&s.BookingID, &s.Title, &s.UserName, &s.StartTime, &s.EndTime); err != nil {
			return nil, err
		}
		slots = append(slots, s)
	}
	return slots, rows.Err()
}
