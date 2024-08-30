package service

import (
	"context"
	"errors"

	"github.com/alinurmyrzakhanov/conference-room-booking/internal/models"
	"github.com/alinurmyrzakhanov/conference-room-booking/internal/repository/postgres"
)

type ReservationService struct {
	repo *postgres.ReservationRepository
}

func NewReservationRepository(repo *postgres.ReservationRepository) *ReservationService {
	return &ReservationService{repo: repo}
}

func (s *ReservationService) CreateReservation(ctx context.Context, reservation models.Reservation) error {
	overlap, err := s.repo.CheckOverlap(ctx, reservation.RoomID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		return err
	}
	if overlap {
		return errors.New("время бронирования занято")
	}
	return s.repo.Create(ctx, reservation)
}
func (s *ReservationService) GetReservations(ctx context.Context, roomID string) ([]models.Reservation, error) {
	return s.repo.GetByRoomId(ctx, roomID)
}
