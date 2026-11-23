package service

import (
	"context"
	"errors"
	"os"
	"strconv"
	"time"

	"university-room-booking/internal/dto"
	"university-room-booking/internal/middleware"
	"university-room-booking/internal/model"
	"university-room-booking/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidCredentials = errors.New("email atau password salah")

type AuthService struct {
	userRepo *repository.UserRepository
}

func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (string, *model.User, error) {
	user, err := s.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}

	expHours, err := strconv.Atoi(os.Getenv("JWT_EXPIRES_HOURS"))
	if err != nil || expHours <= 0 {
		expHours = 24
	}

	claims := middleware.JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", nil, err
	}

	return signedToken, user, nil
}

func (s *AuthService) GetProfile(ctx context.Context, userID uint64) (*model.User, error) {
	return s.userRepo.FindByID(ctx, userID)
}
