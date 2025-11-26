# ---- Etapa de compilación ----
FROM golang:1.22-alpine AS builder

# Instalar certificados (git no es necesario si go.mod está completo)
RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copiar dependencias y descargarlas
COPY go.mod go.sum ./
RUN go mod download

# Copiar el código fuente
COPY . .

# Compilar el binario optimizado
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o webterm ./cmd/webterm

# ---- Etapa final (runtime) ----
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copiar binario y archivos estáticos
COPY --from=builder /app/webterm .
COPY --from=builder /app/static ./static

EXPOSE 8080
ENV PORT=8080

ENTRYPOINT ["./webterm"]