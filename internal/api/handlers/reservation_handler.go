package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/alinurmyrzakhanov/conference-room-booking/internal/models"
	"github.com/alinurmyrzakhanov/conference-room-booking/internal/service"
	"github.com/go-chi/chi/v5"
)

type ReservationHandler struct {
	service *service.ReservationService
}

func NewReservationHandler(service *service.ReservationService) *ReservationHandler {
	return &ReservationHandler{service: service}
}

func (h *ReservationHandler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RoomID    string    `json:"room_id"`
		StartTime time.Time `json:"start_time"`
		EndTime   time.Time `json:"end_time"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	reservation := models.Reservation{
		RoomID:    req.RoomID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	}

	if err := h.service.CreateReservation(r.Context(), reservation); err != nil {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ReservationHandler) GetReservations(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_id")
	reservations, err := h.service.GetReservations(r.Context(), roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(reservations)
}
