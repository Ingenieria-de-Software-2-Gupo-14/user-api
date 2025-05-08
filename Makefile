# Ejecutar la API localmente
run:
	go run ./cmd/api

migrate:
	go run ./cmd/migrations

test:
	go test ./...

# Construir la imagen Docker
docker-build:
	docker build -t api-image .

# Ejecutar la API usando Docker Compose
docker-run:
	docker-compose up --build

swaggo:
	swag init -g ./cmd/api/main.go -o ./docs

# Limpiar archivos generados
clean:
	rm -f main
	docker-compose down
