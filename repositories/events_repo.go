package repositories

import (
	"errors"
	"sync"
	"time"

	"dailyassignment/models"
)

type Repository struct {
	mu       sync.RWMutex
	Events   map[string]*models.Event
	Bookings map[string]*models.Booking
}

func NewRepository() *Repository {
	return &Repository{
		Events: map[string]*models.Event{
			"1": {ID: "1", Title: "Go Seminar", TotalQuota: 10, AvailableQuota: 10},
			"2": {ID: "2", Title: "MVC Workshop", TotalQuota: 5, AvailableQuota: 5},
		},
		Bookings: make(map[string]*models.Booking),
	}
}

func (r *Repository) GetEvents() []*models.Event {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events := make([]*models.Event, 0, len(r.Events))
	for _, e := range r.Events {
		events = append(events, e)
	}
	return events
}

func (r *Repository) BookEvent(userID, eventID string) (*models.Booking, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, b := range r.Bookings {
		if b.UserID == userID && b.EventID == eventID {
			return nil, errors.New("duplicate booking")
		}
	}

	event, exists := r.Events[eventID]
	if !exists {
		return nil, errors.New("event not found")
	}

	if event.AvailableQuota <= 0 {
		return nil, errors.New("quota exceeded")
	}

	event.AvailableQuota--

	bookingID := "b_" + time.Now().Format("20060102150405.000000")
	b := &models.Booking{
		ID:        bookingID,
		EventID:   eventID,
		UserID:    userID,
		CreatedAt: time.Now(),
	}
	r.Bookings[bookingID] = b

	return b, nil
}

func (r *Repository) GetUserBookings(userID string) []*models.Booking {
	r.mu.RLock()
	defer r.mu.RUnlock()

	userBookings := []*models.Booking{}
	for _, b := range r.Bookings {
		if b.UserID == userID {
			userBookings = append(userBookings, b)
		}
	}
	return userBookings
}
