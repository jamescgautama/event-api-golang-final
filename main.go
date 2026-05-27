package main

import (
	"fmt"
	"net/http"

	"dailyassignment/controllers"
	"dailyassignment/middleware"
	"dailyassignment/repositories"
	"dailyassignment/worker"
)

func main() {
	repo := repositories.NewRepository()
	wp := worker.NewPool(3, 10)
	
	authCtrl := &controllers.AuthController{}
	eventCtrl := &controllers.EventController{Repo: repo}
	bookingCtrl := &controllers.BookingController{Repo: repo, WorkerPool: wp}

	mux := http.NewServeMux()

	mux.HandleFunc("/api/login", middleware.LoginRateLimit(authCtrl.Login))
	mux.HandleFunc("/api/events", eventCtrl.ListEvents)
	mux.HandleFunc("/api/bookings", middleware.AuthRequired(middleware.BookingRateLimit(bookingCtrl.BookTicket)))
	mux.HandleFunc("/api/bookings/history", middleware.AuthRequired(bookingCtrl.BookingHistory))

	fmt.Println("Server listening on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		fmt.Printf("Server failed to start: %v\n", err)
	}
}
