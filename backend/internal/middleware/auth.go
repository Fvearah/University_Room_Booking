package middleware

import (
	"net/http"
	"os"
	"strings"

	"university-room-booking/internal/dto"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint64 `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "UNAUTHORIZED", Message: "Token autentikasi tidak ditemukan"},
			})
			return
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &JWTClaims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		}, jwt.WithValidMethods([]string{"HS256"}))

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "INVALID_TOKEN", Message: "Token tidak valid atau sudah kedaluwarsa"},
			})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("user_role", claims.Role)
		c.Next()
	}
}

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get("user_role")
		if !exists {
			c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
				Error: dto.ErrorBody{Code: "FORBIDDEN", Message: "Akses ditolak"},
			})
			return
		}

		userRole := roleVal.(string)
		for _, allowed := range allowedRoles {
			if userRole == allowed {
				c.Next()
				return
			}
		}

		c.AbortWithStatusJSON(http.StatusForbidden, dto.ErrorResponse{
			Error: dto.ErrorBody{Code: "FORBIDDEN", Message: "Anda tidak memiliki hak akses untuk endpoint ini"},
		})
	}
}