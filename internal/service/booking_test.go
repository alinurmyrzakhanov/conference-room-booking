package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/alinurmyrzakhanov/conference-room-booking/internal/models"
)

// MockRepository — это макет репозитория для тестирования.
type MockRepository struct {
	checkOverlapFunc      func(ctx context.Context, roomID string, startTime, endTime time.Time) (bool, error)
	createReservationFunc func(ctx context.Context, reservation models.Reservation) error
	getReservationsFunc   func(ctx context.Context, roomID string) ([]models.Reservation, error)
}

func (m *MockRepository) CheckOverlap(ctx context.Context, roomID string, startTime, endTime time.Time) (bool, error) {
	return m.checkOverlapFunc(ctx, roomID, startTime, endTime)
}

func (m *MockRepository) CreateReservation(ctx context.Context, reservation models.Reservation) error {
	return m.createReservationFunc(ctx, reservation)
}

func (m *MockRepository) GetReservations(ctx context.Context, roomID string) ([]models.Reservation, error) {
	return m.getReservationsFunc(ctx, roomID)
}

func TestCreateReservation(t *testing.T) {
	tests := []struct {
		name           string
		reservation    models.Reservation
		mockOverlap    bool
		mockOverlapErr error
		mockCreateErr  error
		expectedErr    error
	}{
		{
			name: "Successful reservation",
			reservation: models.Reservation{
				RoomID:    "room1",
				StartTime: time.Now(),
				EndTime:   time.Now().Add(time.Hour),
			},
			mockOverlap:   false,
			mockCreateErr: nil,
			expectedErr:   nil,
		},
		{
			name: "Overlapping reservation",
			reservation: models.Reservation{
				RoomID:    "room1",
				StartTime: time.Now(),
				EndTime:   time.Now().Add(time.Hour),
			},
			mockOverlap:   true,
			mockCreateErr: nil,
			expectedErr:   errors.New("reservation overlaps with existing booking"),
		},
		{
			name: "Database error on overlap check",
			reservation: models.Reservation{
				RoomID:    "room1",
				StartTime: time.Now(),
				EndTime:   time.Now().Add(time.Hour),
			},
			mockOverlap:    false,
			mockOverlapErr: errors.New("database error"),
			expectedErr:    errors.New("database error"),
		},
		{
			name: "Database error on create",
			reservation: models.Reservation{
				RoomID:    "room1",
				StartTime: time.Now(),
				EndTime:   time.Now().Add(time.Hour),
			},
			mockOverlap:   false,
			mockCreateErr: errors.New("database error"),
			expectedErr:   errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				checkOverlapFunc: func(ctx context.Context, roomID string, startTime, endTime time.Time) (bool, error) {
					return tt.mockOverlap, tt.mockOverlapErr
				},
				createReservationFunc: func(ctx context.Context, reservation models.Reservation) error {
					return tt.mockCreateErr
				},
			}

			service := NewBookingService(mockRepo)

			err := service.CreateReservation(context.Background(), tt.reservation)

			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) || (err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("CreateReservation() error = %v, expectedErr %v", err, tt.expectedErr)
			}
		})
	}
}

func TestGetReservations(t *testing.T) {
	tests := []struct {
		name             string
		roomID           string
		mockReservations []models.Reservation
		mockErr          error
		expectedErr      error
	}{
		{
			name:   "Successful retrieval",
			roomID: "room1",
			mockReservations: []models.Reservation{
				{ID: 1, RoomID: "room1", StartTime: time.Now(), EndTime: time.Now().Add(time.Hour)},
				{ID: 2, RoomID: "room1", StartTime: time.Now().Add(2 * time.Hour), EndTime: time.Now().Add(3 * time.Hour)},
			},
			mockErr:     nil,
			expectedErr: nil,
		},
		{
			name:             "Database error",
			roomID:           "room1",
			mockReservations: nil,
			mockErr:          errors.New("database error"),
			expectedErr:      errors.New("database error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockRepository{
				getReservationsFunc: func(ctx context.Context, roomID string) ([]models.Reservation, error) {
					return tt.mockReservations, tt.mockErr
				},
			}

			service := NewBookingService(mockRepo)

			reservations, err := service.GetReservations(context.Background(), tt.roomID)

			if (err != nil && tt.expectedErr == nil) || (err == nil && tt.expectedErr != nil) || (err != nil && tt.expectedErr != nil && err.Error() != tt.expectedErr.Error()) {
				t.Errorf("GetReservations() error = %v, expectedErr %v", err, tt.expectedErr)
			}

			if tt.mockErr == nil && len(reservations) != len(tt.mockReservations) {
				t.Errorf("GetReservations() returned %d reservations, expected %d", len(reservations), len(tt.mockReservations))
			}
		})
	}
}
