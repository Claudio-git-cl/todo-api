FROM golang:1.20-alpine AS builder

# Instalar Git
RUN apk add --no-cache git

WORKDIR /app

# Copiar go.mod y go.sum primero
COPY go.mod go.sum ./

# Copiar todo el código fuente
COPY . .

# Compilar la aplicación
RUN go build -o api ./cmd/api

# Imagen final
FROM alpine:3.17

WORKDIR /app

# Copiar el binario compilado
COPY --from=builder /app/api .
COPY --from=builder /app/migrations ./migrations

# Exponer el puerto
EXPOSE 8080

# Ejecutar la aplicación
CMD ["/app/api"]