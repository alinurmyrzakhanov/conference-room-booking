package handlers

import (
	"github.com/alinurmyrzakhanov/conference-room-booking/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func SetupRouter(reservationService *service.ReservationService) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	reservationHandler := NewReservationHandler(reservationService)
	r.Post("/reservations", reservationHandler.CreateReservation)
	r.Get("/reservations/{room_id}", reservationHandler.GetReservations)
	return r
}
