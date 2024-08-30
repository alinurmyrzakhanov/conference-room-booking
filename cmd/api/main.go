package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/alinurmyrzakhanov/conference-room-booking/internal/api/handlers"
	"github.com/alinurmyrzakhanov/conference-room-booking/internal/repository/postgres"
	"github.com/alinurmyrzakhanov/conference-room-booking/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/conference_booking?sslmode=disable"
	}
	dbPool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("не получается подклюиться к базе: %v", err)
	}
	defer dbPool.Close()

	reservationRepo := postgres.NewReservationRepository(dbPool)
	reservationService := service.NewReservationRepository(reservationRepo)
	router := handlers.SetupRouter(reservationService)
	log.Print("Запускается сервер на :8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
