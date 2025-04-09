# Usar una imagen base de Go
FROM golang:1.24

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar los archivos del módulo y descargar dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Copiar el archivo .env
COPY .env .env

# Construir el binario
RUN go build -o main ./cmd/api

# Exponer el puerto definido en el archivo .env
ARG PORT=8080
ENV PORT=${PORT}
EXPOSE ${PORT}

# Comando para ejecutar la aplicación
CMD ["./main"]
