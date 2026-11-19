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

type RoomHandler struct {
	svc *service.RoomService
}

func NewRoomHandler(svc *service.RoomService) *RoomHandler {
	return &RoomHandler{svc: svc}
}

func (h *RoomHandler) List(c *gin.Context) {
	rooms, err := h.svc.List(c.Request.Context(), c.GetString("user_role"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "INTERNAL_ERROR", Message: "Gagal mengambil daftar ruangan"},
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": rooms})
}

func (h *RoomHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "ID ruangan tidak valid"},
		})
		return
	}

	room, err := h.svc.Detail(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, repository.ErrRoomNotFound) {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "NOT_FOUND", Message: "Ruangan tidak ditemukan"},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "INTERNAL_ERROR", Message: "Gagal mengambil detail ruangan"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": room})
}

func (h *RoomHandler) Create(c *gin.Context) {
	var req dto.CreateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "Input tidak valid: " + err.Error()},
		})
		return
	}

	room, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrDuplicateRoomCode):
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "DUPLICATE_CODE", Message: err.Error()},
			})
		case errors.Is(err, service.ErrInvalidRoomInput):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()},
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "INTERNAL_ERROR", Message: "Gagal menambahkan ruangan"},
			})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": room})
}

func (h *RoomHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "ID ruangan tidak valid"},
		})
		return
	}

	var req dto.UpdateRoomRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "Input tidak valid: " + err.Error()},
		})
		return
	}

	room, err := h.svc.Update(c.Request.Context(), id, req)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRoomNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "NOT_FOUND", Message: "Ruangan tidak ditemukan"},
			})
		case errors.Is(err, repository.ErrDuplicateRoomCode):
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "DUPLICATE_CODE", Message: err.Error()},
			})
		case errors.Is(err, service.ErrInvalidRoomInput):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "VALIDATION_ERROR", Message: err.Error()},
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "INTERNAL_ERROR", Message: "Gagal memperbarui ruangan"},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": room})
}

func (h *RoomHandler) Availability(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "ID ruangan tidak valid"},
		})
		return
	}

	res, err := h.svc.GetAvailability(c.Request.Context(), id, c.Query("date"))
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrRoomNotFound):
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "NOT_FOUND", Message: "Ruangan tidak ditemukan"},
			})
		case errors.Is(err, service.ErrInvalidDate):
			c.JSON(http.StatusBadRequest, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "Format tanggal harus YYYY-MM-DD"},
			})
		default:
			c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "INTERNAL_ERROR", Message: "Gagal mengecek ketersediaan ruangan"},
			})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": res})
}
