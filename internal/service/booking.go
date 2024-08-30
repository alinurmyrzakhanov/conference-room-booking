package service

import (
	"context"
	"errors"
	"time"

	"github.com/alinurmyrzakhanov/conference-room-booking/internal/models"
)

type Repository interface {
	CreateReservation(ctx context.Context, reservation models.Reservation) error
	GetReservations(ctx context.Context, roomID string) ([]models.Reservation, error)
	CheckOverlap(ctx context.Context, roomID string, startTime, endTime time.Time) (bool, error)
}

type BookingService struct {
	repo Repository
}

func NewBookingService(repo Repository) *BookingService {
	return &BookingService{repo: repo}
}

func (s *BookingService) CreateReservation(ctx context.Context, reservation models.Reservation) error {
	// Проверка на пересечение бронирований
	overlap, err := s.repo.CheckOverlap(ctx, reservation.RoomID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		return err
	}
	if overlap {
		return errors.New("reservation overlaps with existing booking")
	}

	// Создание бронирования
	return s.repo.CreateReservation(ctx, reservation)
}

func (s *BookingService) GetReservations(ctx context.Context, roomID string) ([]models.Reservation, error) {
	return s.repo.GetReservations(ctx, roomID)
}
