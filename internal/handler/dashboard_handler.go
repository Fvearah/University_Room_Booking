package handler

import (
	"net/http"

	"university-room-booking/internal/dto"
	"university-room-booking/internal/service"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler(svc *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{svc: svc}
}

func (h *DashboardHandler) Summary(c *gin.Context) {
	summary, err := h.svc.GetSummary(c.Request.Context(), c.GetUint64("user_id"), c.GetString("user_role"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "INTERNAL_ERROR", Message: "Gagal memuat ringkasan dashboard"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": summary})
}
