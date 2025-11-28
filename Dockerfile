
FROM golang:1.25-alpine AS builder


RUN apk add --no-cache ca-certificates

WORKDIR /app


COPY go.mod go.sum ./
RUN go mod download


COPY . .


RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o webterm ./cmd/webterm


FROM alpine:3.18

RUN apk add --no-cache ca-certificates

WORKDIR /app


COPY --from=builder /app/webterm .
COPY --from=builder /app/static ./static

EXPOSE 8080
ENV PORT=8080

ENTRYPOINT ["./webterm"]