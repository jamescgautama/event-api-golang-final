package controllers

import (
	"encoding/json"
	"net/http"

	"dailyassignment/middleware"
	"dailyassignment/repositories"
	"dailyassignment/worker"
)

type BookingController struct {
	Repo *repositories.Repository
	WorkerPool *worker.Pool
}

type BookingRequest struct {
	EventID string `json:"event_id"`
}

func (c *BookingController) BookTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)

	var req BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	booking, err := c.Repo.BookEvent(userID, req.EventID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Enqueue job for background processing
	if c.WorkerPool != nil {
		c.WorkerPool.Enqueue(worker.Job{
			UserID:  userID,
			EventID: req.EventID,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

func (c *BookingController) BookingHistory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(string)
	bookings := c.Repo.GetUserBookings(userID)
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}
