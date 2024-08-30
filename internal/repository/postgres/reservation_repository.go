package postgres

import (
	"context"
	"time"

	"github.com/alinurmyrzakhanov/conference-room-booking/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ReservationRepository struct {
	db *pgxpool.Pool
}

func NewReservationRepository(db *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (r *ReservationRepository) Create(ctx context.Context, reservation models.Reservation) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err := r.db.Exec(ctx,
		"INSERT INTO reservations (room_id, start_time, end_time) VALUES ($1, $2, $3)",
		reservation.RoomID, reservation.StartTime, reservation.EndTime)
	return err
}

func (r *ReservationRepository) GetByRoomId(ctx context.Context, roomID string) ([]models.Reservation, error) {
	rows, err := r.db.Query(ctx,
		"SELECT id, room_id, start_time, end_time FROM reservations WHERE room_id = $1", roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var reservations []models.Reservation
	for rows.Next() {
		var res models.Reservation
		err := rows.Scan(&res.ID, &res.RoomID, &res.StartTime, &res.EndTime)
		if err != nil {
			return nil, err
		}
		reservations = append(reservations, res)
	}
	return reservations, nil

}

func (r *ReservationRepository) CheckOverlap(ctx context.Context, roomID string, startTime, endTime time.Time) (bool, error) {
	var count int
	err := r.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM reservations
	WHERE room_id = $1 AND ((start_time <= $2 AND end_time > $2) OR (start_time < $3 AND end_time >= $3) OR (start_time >= $2 AND end_time <= $3))`,
		roomID, startTime, endTime).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, err

}
