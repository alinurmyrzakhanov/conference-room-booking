package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alinurmyrzakhanov/conference-room-booking/internal/api/handlers"
	"github.com/alinurmyrzakhanov/conference-room-booking/internal/config"
	repository "github.com/alinurmyrzakhanov/conference-room-booking/internal/repository/postgres"
	"github.com/alinurmyrzakhanov/conference-room-booking/internal/service"
	"github.com/jackc/pgx/v4/pgxpool"
)

func main() {
	// Загрузка конфигурации
	cfg := config.NewConfig()

	// Подключение к базе данных
	dbPool, err := pgxpool.Connect(context.Background(), cfg.DBUrl)
	if err != nil {
		log.Fatalf("Не удается подключиться к базе: %v", err)
	}
	defer dbPool.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := dbPool.Ping(ctx); err != nil {
		log.Fatalf("пинг базы данных не прошел: %v", err)
	}

	// Инициализация репозитория и сервиса
	reservationRepo := repository.NewReservationRepository(dbPool)
	reservationService := service.NewReservationService(reservationRepo)
	// Настройка маршрутизации
	router := handlers.SetupRouter(reservationService)

	// Запуск сервера
	serverAddr := fmt.Sprintf(":%d", cfg.ServerPort)
	log.Printf("Запускается сервер на %s", serverAddr)
	if err := http.ListenAndServe(serverAddr, router); err != nil {
		log.Fatalf("Не удалось запустить сервер: %v", err)
	}
}
