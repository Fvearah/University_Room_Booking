package handler

import (
	"errors"
	"net/http"
	"strconv"

	"university-room-booking/internal/dto"
	"university-room-booking/internal/repository"
	"university-room-booking/internal/service"

	"github.com/gin-gonic/gin"
)

type BookingHandler struct {
	svc *service.BookingService
}

func NewBookingHandler(svc *service.BookingService) *BookingHandler {
	return &BookingHandler{svc: svc}
}

func (h *BookingHandler) Create(c *gin.Context) {
	var req dto.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "Input tidak valid: " + err.Error()},
		})
		return
	}

	userID := c.GetUint64("user_id")
	booking, err := h.svc.Create(c.Request.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRoomConflict):
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error: dto.ErrorBody{
					Code:    "ROOM_UNAVAILABLE",
					Message: "Ruangan sudah dipesan pada waktu tersebut. Silakan pilih jadwal lain.",
				},
			})
		case errors.Is(err, repository.ErrRoomNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "ROOM_NOT_FOUND", Message: "Ruangan tidak ditemukan"},
			})
		case errors.Is(err, repository.ErrRoomInactive), errors.Is(err, service.ErrInvalidDate), errors.Is(err, service.ErrInvalidTime):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()},
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "INTERNAL_ERROR", Message: "Terjadi kesalahan pada server"},
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": booking})
}

func (h *BookingHandler) List(c *gin.Context) {
	userID := c.GetUint64("user_id")
	role := c.GetString("user_role")

	bookings, err := h.svc.List(
		c.Request.Context(),
		userID,
		role,
		c.Query("date"),
		c.Query("room_id"),
		c.Query("status"),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "INTERNAL_ERROR", Message: "Gagal mengambil daftar booking"},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": bookings})
}

func (h *BookingHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "ID booking tidak valid"},
		})
		return
	}

	booking, err := h.svc.GetDetail(c.Request.Context(), id, c.GetUint64("user_id"), c.GetString("user_role"))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrBookingNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "NOT_FOUND", Message: "Data booking tidak ditemukan"},
			})
		case errors.Is(err, service.ErrForbiddenAccess):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "FORBIDDEN", Message: err.Error()},
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "INTERNAL_ERROR", Message: "Gagal mengambil detail booking"},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": booking})
}

func (h *BookingHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "ID booking tidak valid"},
		})
		return
	}

	err = h.svc.Cancel(c.Request.Context(), id, c.GetUint64("user_id"))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrBookingNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "NOT_FOUND", Message: "Data booking tidak ditemukan"},
			})
		case errors.Is(err, service.ErrForbiddenAccess):
			c.JSON(http.StatusForbidden, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "FORBIDDEN", Message: err.Error()},
			})
		case errors.Is(err, service.ErrCannotCancel):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "CANNOT_CANCEL", Message: err.Error()},
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "INTERNAL_ERROR", Message: "Gagal membatalkan booking"},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Booking berhasil dibatalkan"})
}
