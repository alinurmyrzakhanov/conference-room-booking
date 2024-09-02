package repository

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/alinurmyrzakhanov/conference-room-booking/internal/models"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ReservationRepositoryInterface interface {
	Create(ctx context.Context, reservation models.Reservation) error
	GetByRoomID(ctx context.Context, roomID string) ([]models.Reservation, error)
	CheckOverlap(ctx context.Context, roomID string, startTime, endTime time.Time) (bool, error)
}

type ReservationRepository struct {
	db *pgxpool.Pool
}

func NewReservationRepository(db *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{db: db}
}

// Create добавляет новое бронирование в базу данных
func (r *ReservationRepository) Create(ctx context.Context, reservation models.Reservation) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query, args, err := squirrel.Insert("reservations").
		Columns("room_id", "start_time", "end_time").
		Values(reservation.RoomID, reservation.StartTime, reservation.EndTime).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = r.db.Exec(ctx, query, args...)
	return err
}

type ReservationData struct {
	RoomID    string    `json:"room_id"`
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// GetByRoomID возвращает все бронирования для указанной комнаты
func (r *ReservationRepository) GetByRoomID(ctx context.Context, roomID string) ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query, args, err := squirrel.Select("room_id", "start_time", "end_time").
		From("reservations").
		Where(squirrel.Eq{"room_id": roomID}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return nil, err
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservationsData []ReservationData

	for rows.Next() {
		var resData ReservationData
		if err := rows.Scan(&resData.RoomID, &resData.StartTime, &resData.EndTime); err != nil {
			return nil, err
		}
		reservationsData = append(reservationsData, resData)
	}
	var reservations []models.Reservation
	for _, data := range reservationsData {
		res := models.Reservation{
			RoomID:    data.RoomID,
			StartTime: data.StartTime,
			EndTime:   data.EndTime,
		}
		reservations = append(reservations, res)
	}
	return reservations, nil
}

// CheckOverlap проверяет, пересекается ли новое бронирование с существующими
func (r *ReservationRepository) CheckOverlap(ctx context.Context, roomID string, startTime, endTime time.Time) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query, args, err := squirrel.Select("COUNT(*)").
		From("reservations").
		Where(squirrel.Eq{"room_id": roomID}).
		Where(squirrel.Or{
			squirrel.And{
				squirrel.LtOrEq{"start_time": startTime},
				squirrel.Gt{"end_time": startTime},
			},
			squirrel.And{
				squirrel.Lt{"start_time": endTime},
				squirrel.GtOrEq{"end_time": endTime},
			},
			squirrel.And{
				squirrel.GtOrEq{"start_time": startTime},
				squirrel.LtOrEq{"end_time": endTime},
			},
		}).
		PlaceholderFormat(squirrel.Dollar).
		ToSql()

	if err != nil {
		return false, err
	}

	var count int
	err = r.db.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
