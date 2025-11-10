# ---- Etapa de compilación ----
FROM golang:1.25-alpine AS builder

# Instalar git y certificados (Go los necesita)
RUN apk add --no-cache git ca-certificates

# Crear directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias y descargarlas
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código del proyecto
COPY . .

# Compilar el binario
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o webterm ./cmd/webterm

# ---- Etapa final (runtime) ----
FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /root/

# Copiar el binario y los archivos estáticos
COPY --from=builder /app/webterm .
COPY --from=builder /app/static ./static

# Puerto que usa la app
EXPOSE 8080

# Variable de entorno para el puerto
ENV PORT=8080

# Comando que ejecuta el programa
CMD ["./webterm"]
