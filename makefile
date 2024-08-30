.PHONY: build run test

build:
	go build -o app ./cmd/api

run:
	go run ./cmd/api

test:
	go test -v ./...

.PHONY: docker-build
docker-build:
	docker-compose build

.PHONY: docker-up
docker-up:
	docker-compose up -d

.PHONY: docker-down
docker-down:
	docker-compose down

.PHONY: docker-test
docker-test:
	docker-compose run --rm api go test -v ./...