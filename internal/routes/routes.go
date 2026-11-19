package routes

import (
	"university-room-booking/internal/handler"
	"university-room-booking/internal/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authH *handler.AuthHandler,
	roomH *handler.RoomHandler,
	bookingH *handler.BookingHandler,
	dashboardH *handler.DashboardHandler,
) *gin.Engine {
	r := gin.Default()

	// Middleware CORS untuk integrasi Frontend
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	v1 := r.Group("/api/v1")
	{
		// Endpoint Publik
		v1.POST("/auth/login", authH.Login)

		// Endpoint Privat (Memerlukan Token JWT)
		protected := v1.Group("")
		protected.Use(middleware.AuthRequired())
		{
			protected.GET("/auth/me", authH.Me)
			protected.GET("/dashboard", dashboardH.Summary)

			// Ruangan
			protected.GET("/rooms", roomH.List)
			protected.GET("/rooms/:id", roomH.Detail)
			protected.GET("/rooms/:id/availability", roomH.Availability)

			adminRooms := protected.Group("/rooms")
			adminRooms.Use(middleware.RequireRole("ADMIN"))
			{
				adminRooms.POST("", roomH.Create)
				adminRooms.PATCH("/:id", roomH.Update)
			}

			// Booking
			protected.GET("/bookings", bookingH.List)
			protected.GET("/bookings/:id", bookingH.Detail)

			userBookings := protected.Group("/bookings")
			userBookings.Use(middleware.RequireRole("USER"))
			{
				userBookings.POST("", bookingH.Create)
				userBookings.PATCH("/:id/cancel", bookingH.Cancel)
			}
		}
	}

	return r
}
