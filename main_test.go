package main_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"context"

	"dailyassignment/controllers"
	"dailyassignment/middleware"
	"dailyassignment/repositories"
)

func TestBookingConcurrency(t *testing.T) {
	repo := repositories.NewRepository()
	// set initial quota for an event to 100 for testing
	repo.Events["1"].TotalQuota = 100
	repo.Events["1"].AvailableQuota = 100

	bookingCtrl := &controllers.BookingController{Repo: repo}

	var wg sync.WaitGroup
	// we will attempt 200 concurrent bookings for event "1", only 100 should succeed, rest gets fucked
	numRequests := 200
	successCount := 0
	var mu sync.Mutex

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(userID string) {
			defer wg.Done()

			reqBody := `{"event_id": "1"}`
			req, _ := http.NewRequest("POST", "/api/bookings", strings.NewReader(reqBody))
			// bypass standard auth middleware by injecting user ID directly to request context
			ctx := context.WithValue(req.Context(), middleware.UserIDKey, userID)
			req = req.WithContext(ctx)
			
			w := httptest.NewRecorder()
			bookingCtrl.BookTicket(w, req)

			if w.Code == http.StatusOK {
				mu.Lock()
				successCount++
				mu.Unlock()
			}
		}(fmt.Sprintf("user_%d", i))
	}

	wg.Wait()

	if successCount != 100 {
		t.Errorf("Expected exactly 100 successful bookings, got %d", successCount)
	}

	available := repo.Events["1"].AvailableQuota
	if available != 0 {
		t.Errorf("Expected available quota to be 0, got %d", available)
	}
}
