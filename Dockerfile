
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache tzdata ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

CMD ["go", "run", "main.go"]