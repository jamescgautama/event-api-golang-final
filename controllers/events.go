package controllers

import (
	"encoding/json"
	"net/http"

	"dailyassignment/repositories"
)

type EventController struct {
	Repo *repositories.Repository
}

func (c *EventController) ListEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	events := c.Repo.GetEvents()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(events)
}
