package tests

import (
	"context"
	"log"
	"os"
	"testing"
	"time"

	"github.com/alinurmyrzakhanov/conference-room-booking/internal/models"
	repository "github.com/alinurmyrzakhanov/conference-room-booking/internal/repository/postgres"
	"github.com/alinurmyrzakhanov/conference-room-booking/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
)

func TestReservationService(t *testing.T) {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://user:password@localhost:5432/conference_booking?sslmode=disable"
	}
	dbPool, err := pgxpool.Connect(context.Background(), dbURL)
	if err != nil {
		log.Fatalf("не получается подклюиться к базе: %v", err)
	}
	defer dbPool.Close()

	_, err = dbPool.Exec(context.Background(), "TRUNCATE TABLE reservations")
	if err != nil {
		t.Fatalf("ошибка обрезания таблицы: %v", err)
	}

	repo := repository.NewReservationRepository(dbPool)
	svc := service.NewReservationService(repo)

	t.Run("CreateReservation", func(t *testing.T) {
		reservation := models.Reservation{
			RoomID:    "room1",
			StartTime: time.Now(),
			EndTime:   time.Now().Add(time.Hour),
		}

		err := svc.CreateReservation(context.Background(), reservation)
		if err != nil {
			t.Errorf("ошибка создания брони: %v", err)
		}
	})

	t.Run("CreateOverlappingReservation", func(t *testing.T) {
		reservation := models.Reservation{
			RoomID:    "room1",
			StartTime: time.Now().Add(30 * time.Minute),
			EndTime:   time.Now().Add(90 * time.Minute),
		}

		err := svc.CreateReservation(context.Background(), reservation)
		if err == nil {
			t.Error("ожидаемая ошибка при наличии брони")
		}
	})

	t.Run("GetReservations", func(t *testing.T) {
		reservations, err := svc.GetReservations(context.Background(), "room1")
		if err != nil {
			t.Errorf("ошибка получения брони: %v", err)
		}
		if len(reservations) != 1 {
			t.Errorf("ожидалась одна бронь получили %d", len(reservations))
		}
	})
}
