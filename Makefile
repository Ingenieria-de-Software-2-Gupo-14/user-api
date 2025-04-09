# Ejecutar la API localmente
run:
	go run ./cmd/api

test:
	go test ./...

# Construir la imagen Docker
docker-build:
	docker build -t api-image .

# Ejecutar la API usando Docker Compose
docker-run:
	docker-compose up --build

# Limpiar archivos generados
clean:
	rm -f main
	docker-compose down
