package service

import (
	"context"
	"errors"

	"github.com/alinurmyrzakhanov/conference-room-booking/internal/models"
	repository "github.com/alinurmyrzakhanov/conference-room-booking/internal/repository/postgres"
)

type ReservationServiceInterface interface {
	CreateReservation(ctx context.Context, reservation models.Reservation) error
	GetReservations(ctx context.Context, roomID string) ([]models.Reservation, error)
}

type ReservationService struct {
	repo repository.ReservationRepositoryInterface
}

func NewReservationService(repo repository.ReservationRepositoryInterface) *ReservationService {
	return &ReservationService{repo: repo}
}

// CreateReservation создает новое бронирование, если оно не пересекается с существующими
func (s *ReservationService) CreateReservation(ctx context.Context, reservation models.Reservation) error {
	overlap, err := s.repo.CheckOverlap(ctx, reservation.RoomID, reservation.StartTime, reservation.EndTime)
	if err != nil {
		return err
	}
	if overlap {
		return errors.New("время занято")
	}
	return s.repo.Create(ctx, reservation)
}

// GetReservations возвращает все бронирования для указанной комнаты
func (s *ReservationService) GetReservations(ctx context.Context, roomID string) ([]models.Reservation, error) {
	return s.repo.GetByRoomID(ctx, roomID)
}
