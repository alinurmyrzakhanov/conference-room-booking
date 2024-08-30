# Система бронирования конференц-залов

Это мой REST API для системы бронирования конференц-залов

## Запуск приложения

1. Скачайте себе репо:
git clone https://github.com/alinur-rama/conference-room-booking.git
cd conference-room-booking

2. Запуск приложения:
у вас должен быть установлен docker desktop
docker-compose up --build
API будет доступно по адресу `http://localhost:8080`.

## API Endpoints

- `curl -X POST http://localhost:8080/reservations      -H "Content-Type: application/json"      -d '{
           "room_id": "room1",
           "start_time": "2024-08-30T10:00:00Z",
           "end_time": "2024-08-30T11:00:00Z"
         }'
`: Создание нового бронирования
- `curl -X GET http://localhost:8080/reservations/room1`: Получение всех бронирований для комнаты
- можете повторить верхний пост чтобы увидеть что время занято

## Запуск тестов
docker-compose run --rm tests

## Структура проекта

- `cmd/api`: Основная точка входа приложения
- `internal/api`: Обработчики API и маршрутизатор
