package handler

import (
	"errors"
	"net/http"

	"university-room-booking/internal/dto"
	"university-room-booking/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler(svc *service.AuthService) *AuthHandler {
	return &AuthHandler{svc: svc}
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "BAD_REQUEST", Message: "Format email atau password tidak valid"},
		})
		return
	}

	token, user, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			c.JSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "INVALID_CREDENTIALS", Message: err.Error()},
			})
			return
		}
		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "INTERNAL_ERROR", Message: "Gagal memproses login"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"token": token,
			"user":  user,
		},
	})
}

func (h *AuthHandler) Me(c *gin.Context) {
	userID := c.GetUint64("user_id")
	user, err := h.svc.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "NOT_FOUND", Message: "Pengguna tidak ditemukan"},
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": user})
}
