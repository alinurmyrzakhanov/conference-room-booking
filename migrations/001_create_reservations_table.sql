CREATE TABLE IF NOT EXISTS reservations(
    id SERIAL PRIMARY KEY,
    room_id VARCHAR(255) NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL
);

CREATE INDEX idx_room_id_start_time_end_time ON reservations(room_id, start_time, end_time)