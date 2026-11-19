package service

import (
	"context"
	"time"

	"university-room-booking/internal/repository"
)

type DashboardService struct {
	repo *repository.DashboardRepository
}

func NewDashboardService(repo *repository.DashboardRepository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) GetSummary(ctx context.Context, userID uint64, role string) (interface{}, error) {
	if role == "ADMIN" {
		return s.repo.GetAdminStats(ctx)
	}
	now := time.Now().In(time.Local)
	return s.repo.GetUserStats(ctx, userID, now.Format("2006-01-02"), now.Format("15:04:05"))
}
