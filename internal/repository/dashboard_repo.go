package repository

import (
	"context"
	"database/sql"
)

type AdminDashboardStats struct {
	TotalRooms        int `json:"total_rooms"`
	ActiveRooms       int `json:"active_rooms"`
	InactiveRooms     int `json:"inactive_rooms"`
	TotalBookings     int `json:"total_bookings"`
	ConfirmedBookings int `json:"confirmed_bookings"`
	CancelledBookings int `json:"cancelled_bookings"`
}

type UserDashboardStats struct {
	TotalBookings     int `json:"total_bookings"`
	UpcomingBookings  int `json:"upcoming_bookings"`
	CompletedBookings int `json:"completed_bookings"`
	CancelledBookings int `json:"cancelled_bookings"`
}

type DashboardRepository struct {
	db *sql.DB
}

func NewDashboardRepository(db *sql.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (r *DashboardRepository) GetAdminStats(ctx context.Context) (*AdminDashboardStats, error) {
	var stats AdminDashboardStats

	roomQuery := `
		SELECT
			COUNT(*) AS total,
			COALESCE(SUM(CASE WHEN status = 'ACTIVE' THEN 1 ELSE 0 END), 0) AS active_count,
			COALESCE(SUM(CASE WHEN status = 'INACTIVE' THEN 1 ELSE 0 END), 0) AS inactive_count
		FROM rooms`
	if err := r.db.QueryRowContext(ctx, roomQuery).Scan(&stats.TotalRooms, &stats.ActiveRooms, &stats.InactiveRooms); err != nil {
		return nil, err
	}

	bookingQuery := `
		SELECT
			COUNT(*) AS total,
			COALESCE(SUM(CASE WHEN status = 'CONFIRMED' THEN 1 ELSE 0 END), 0) AS confirmed_count,
			COALESCE(SUM(CASE WHEN status = 'CANCELLED' THEN 1 ELSE 0 END), 0) AS cancelled_count
		FROM bookings`
	if err := r.db.QueryRowContext(ctx, bookingQuery).Scan(&stats.TotalBookings, &stats.ConfirmedBookings, &stats.CancelledBookings); err != nil {
		return nil, err
	}

	return &stats, nil
}

func (r *DashboardRepository) GetUserStats(ctx context.Context, userID uint64, currentDate, currentTime string) (*UserDashboardStats, error) {
	var stats UserDashboardStats

	query := `
		SELECT
			COUNT(*) AS total,
			COALESCE(SUM(CASE
				WHEN status = 'CONFIRMED' AND (booking_date > ? OR (booking_date = ? AND end_time >= ?))
				THEN 1 ELSE 0
			END), 0) AS upcoming_count,
			COALESCE(SUM(CASE
				WHEN status = 'CONFIRMED' AND (booking_date < ? OR (booking_date = ? AND end_time < ?))
				THEN 1 ELSE 0
			END), 0) AS completed_count,
			COALESCE(SUM(CASE WHEN status = 'CANCELLED' THEN 1 ELSE 0 END), 0) AS cancelled_count
		FROM bookings
		WHERE user_id = ?`

	err := r.db.QueryRowContext(
		ctx, query,
		currentDate, currentDate, currentTime,
		currentDate, currentDate, currentTime,
		userID,
	).Scan(&stats.TotalBookings, &stats.UpcomingBookings, &stats.CompletedBookings, &stats.CancelledBookings)
	if err != nil {
		return nil, err
	}

	return &stats, nil
}
