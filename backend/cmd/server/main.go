package main

import (
	"log"
	"os"

	"university-room-booking/internal/config"
	"university-room-booking/internal/handler"
	"university-room-booking/internal/repository"
	"university-room-booking/internal/routes"
	"university-room-booking/internal/service"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

	if os.Getenv("JWT_SECRET") == "" {
		log.Fatal("JWT_SECRET belum diisi di file .env")
	}

	db, err := config.InitDB()
	if err != nil {
		log.Fatalf("Gagal terhubung ke database MySQL: %v", err)
	}
	defer db.Close()

	// Inisialisasi Repository
	userRepo := repository.NewUserRepository(db)
	roomRepo := repository.NewRoomRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	dashboardRepo := repository.NewDashboardRepository(db)

	// Inisialisasi Service
	authSvc := service.NewAuthService(userRepo)
	roomSvc := service.NewRoomService(roomRepo)
	bookingSvc := service.NewBookingService(bookingRepo)
	dashboardSvc := service.NewDashboardService(dashboardRepo)

	// Inisialisasi Handler
	authHandler := handler.NewAuthHandler(authSvc)
	roomHandler := handler.NewRoomHandler(roomSvc)
	bookingHandler := handler.NewBookingHandler(bookingSvc)
	dashboardHandler := handler.NewDashboardHandler(dashboardSvc)

	router := routes.SetupRouter(authHandler, roomHandler, bookingHandler, dashboardHandler)

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server berjalan pada http://localhost:%s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}
